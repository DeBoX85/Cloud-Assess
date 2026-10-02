// Query reproduced from Microsoft Azure Quick Review (MIT licensed). See NOTICE.md and docs/SERVICE_HEALTH.md.
package servicehealth

const Query = `servicehealthresources
| where type =~ 'Microsoft.ResourceHealth/events' and properties.Status == 'Resolved'
  and (properties.EventType == 'ServiceIssue')
| extend
    eventTrackingId  = tostring(id),
    subscriptionId   = tostring(subscriptionId),
    impactStart      = todatetime(tostring(properties.ImpactStartTime)),
    impactMitigation = todatetime(tostring(properties.ImpactMitigationTime))
| where isnotnull(impactMitigation) and impactMitigation >= impactStart
| extend durationHours = max_of(0.0, round((impactMitigation - impactStart) / 1h, 3))
| join kind=inner (
    servicehealthresources
    | where type == 'microsoft.resourcehealth/events/impactedresources'
    | extend eventTrackingId = tostring(split(id, '/impactedResources/')[0])
    | extend p = parse_json(properties)
    | project
        eventTrackingId,
        targetResourceId   = tostring(p.targetResourceId),
        targetResourceType = tostring(p.targetResourceType),
        targetRegion       = tostring(p.targetRegion)
) on $left.eventTrackingId == $right.eventTrackingId
| project subscriptionId, targetRegion, targetResourceId, targetResourceType, durationHours
| summarize totalDurationHours = sum(durationHours)
    by targetResourceId, subscriptionId, targetRegion, targetResourceType
| extend totalDurationHours = min_of(toreal(2160), totalDurationHours)
| extend percentageOfTimeWithoutEvents = max_of(0.0, round((2160 - totalDurationHours) * 100 / 2160, 3))
| project subscriptionId, targetRegion, targetResourceId, targetResourceType, percentageOfTimeWithoutEvents
| union (
    resources
    | project subscriptionId = tostring(subscriptionId), targetRegion = tostring(location),
              targetResourceId = tostring(id), targetResourceType = tostring(type)
    | join kind=leftouter (
        servicehealthresources
        | where type == 'microsoft.resourcehealth/events/impactedresources'
        | extend p = parse_json(properties)
        | project targetResourceId = tostring(p.targetResourceId)
        | distinct targetResourceId
    ) on targetResourceId
    | where isempty(targetResourceId1)
    | project subscriptionId, targetRegion, targetResourceId, targetResourceType,
        percentageOfTimeWithoutEvents = toreal(100)
)
| summarize
    percentageOfTimeWithoutEvents = round(avg(percentageOfTimeWithoutEvents), 2),
    affectedResources             = countif(percentageOfTimeWithoutEvents < 100)
    by subscriptionId, targetRegion, targetResourceType
| join kind=leftouter (
    servicehealthresources
    | where type =~ 'Microsoft.ResourceHealth/events' and properties.Status == 'Resolved'
      and (properties.EventType == 'ServiceIssue' or properties.EventType == 'PlannedMaintenance')
    | extend eventId = tostring(id), subscriptionId = tostring(subscriptionId)
    | join kind=inner (
        servicehealthresources
        | where type == 'microsoft.resourcehealth/events/impactedresources'
        | extend eventId = tostring(split(id, '/impactedResources/')[0])
        | extend p = parse_json(properties)
        | project eventId, targetResourceType = tostring(p.targetResourceType),
                  targetRegion = tostring(p.targetRegion)
    ) on eventId
    | summarize events = dcount(eventId) by subscriptionId, targetRegion, targetResourceType
) on subscriptionId, targetRegion, targetResourceType
| project
    subscriptionId,
    targetRegion,
    targetResourceType,
    percentageOfTimeWithoutEvents,
    events            = coalesce(events, 0),
    affectedResources
| order by percentageOfTimeWithoutEvents asc`
