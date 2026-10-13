package region

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

// Declared from the pinned source contract before reading emitted observations.
const otherSourceIndependentOracles = `{
  "storage-empty": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-null": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-null-items": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-root-array": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-trailing": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-cancelled": {
    "Entries": null,
    "Failed": true,
    "URLs": [],
    "DisplayRegion": false
  },
  "cognitive-empty": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-null": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-null-items": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-root-array": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-trailing": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-cancelled": {
    "Entries": null,
    "Failed": true,
    "URLs": [],
    "DisplayRegion": false
  },
  "sql-db-empty": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-null": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-null-items": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-root-array": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-trailing": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-cancelled": {
    "Entries": null,
    "Failed": true,
    "URLs": [],
    "DisplayRegion": false
  },
  "sql-mi-empty": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-null": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-null-items": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-root-array": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-trailing": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-cancelled": {
    "Entries": null,
    "Failed": true,
    "URLs": [],
    "DisplayRegion": false
  },
  "storage-single": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-empty-list": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-null-list": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-wrong-list": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-wrong-item": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-wrong-name": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-wrong-locations": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-wrong-next": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-missing-locations": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-null-locations": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-empty-locations": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-global-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-display-spaces": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-tabs-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-nbsp-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-unicode-fold-location": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-global-explicit": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-whitespace-name": {
    "Entries": {
      " s ": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-surrogate-name": {
    "Entries": {
      "�": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-single-field-alias": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-last-location-alias": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-last-value-alias": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-duplicate-name": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-paged": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      },
      "second": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-later-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "storage-next-foreign": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://synthetic.invalid/other/next"
    ],
    "DisplayRegion": false
  },
  "storage-next-relative": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-next-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-bounded-cycle": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-single": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-empty-list": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-null-list": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-list": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-item": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-name": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-locations": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-next": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-missing-locations": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-null-locations": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-empty-locations": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-global-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-display-spaces": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-tabs-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-nbsp-excluded": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-unicode-fold-location": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-global-explicit": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-whitespace-name": {
    "Entries": {
      " s ": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-surrogate-name": {
    "Entries": {
      "�": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-single-field-alias": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-last-location-alias": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-last-value-alias": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-duplicate-name": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-paged": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      },
      "second": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-http400": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-http403": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-http404": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-http429": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-http500": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-later-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-next-foreign": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://synthetic.invalid/other/next"
    ],
    "DisplayRegion": false
  },
  "cognitive-next-relative": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-next-malformed": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-bounded-cycle": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "storage-name-tier": {
    "Entries": {
      "standard_lrs": {
        "State": 0,
        "BlockedZones": null
      },
      "premium_lrs": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-location-values-ignored": {
    "Entries": {
      "s": {
        "State": 1,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-quotaid": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-location-priority": {
    "Entries": {
      "s": {
        "State": 1,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-zone-order": {
    "Entries": {
      "s": {
        "State": 3,
        "BlockedZones": [
          "2",
          "1",
          "2"
        ]
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-empty-zone-capability": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-capability-true": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-last-state": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-unicode-restriction": {
    "Entries": {
      "s": {
        "State": 1,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-invalid-zones": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "storage-ignored-kind": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-tier-not-fallback": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-kinds-collapse": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-restrictions-ignored": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-malformed-restrictions-ignored": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-kind": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "cognitive-wrong-resource-type": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-null-roots": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-server-root": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-mi-root": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-editions": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-slos": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-sku": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-sku-tier": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-status": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-wrong-mi-family": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-unknown-value-ignored": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-null-inner": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-null-roots": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-server-root": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-mi-root": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-editions": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-slos": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-sku": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-sku-tier": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-status": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-wrong-mi-family": {
    "Entries": null,
    "Failed": true,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-unknown-value-ignored": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-null-inner": {
    "Entries": {},
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-status-matrix": {
    "Entries": {
      "available": {
        "State": 0,
        "BlockedZones": null
      },
      "default": {
        "State": 0,
        "BlockedZones": null
      },
      "visible": {
        "State": 0,
        "BlockedZones": null
      },
      "unknown": {
        "State": 0,
        "BlockedZones": null
      },
      "missing": {
        "State": 0,
        "BlockedZones": null
      },
      "null": {
        "State": 0,
        "BlockedZones": null
      },
      "padded": {
        "State": 0,
        "BlockedZones": null
      },
      "disabled": {
        "State": 2,
        "BlockedZones": null
      },
      "unicodedisabled": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-name-no-tier": {
    "Entries": {
      " s ": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-last-state": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-parent-status-ignored": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-multi-level": {
    "Entries": {
      "a": {
        "State": 0,
        "BlockedZones": null
      },
      "c": {
        "State": 0,
        "BlockedZones": null
      },
      "b": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-unicode-root-alias": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-last-root-alias": {
    "Entries": {
      "s": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-next-ignored": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-tier-matrix": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      },
      "bc_gen8im": {
        "State": 0,
        "BlockedZones": null
      },
      "ngpgen5_gen5": {
        "State": 0,
        "BlockedZones": null
      },
      "raw_family": {
        "State": 0,
        "BlockedZones": null
      },
      "_gen5": {
        "State": 0,
        "BlockedZones": null
      },
      " padded _ gen5 ": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-status-ignored": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-case-collision": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-name-null": {
    "Entries": {
      "_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-unicode-root-alias": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-next-ignored": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-shared-payload": {
    "Entries": {
      "sql_d": {
        "State": 2,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-shared-payload": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-direct-upper-region": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/EASTUS/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-mi-direct-upper-region": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/EASTUS/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  },
  "sql-db-direct-display-region": {
    "Entries": {
      "s": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/East%20US/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": true
  },
  "sql-mi-direct-display-region": {
    "Entries": {
      "gp_gen5": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/East%20US/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": true
  },
  "storage-paged-last-state": {
    "Entries": {
      "s": {
        "State": 1,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "cognitive-paged-kinds-collapse": {
    "Entries": {
      "s0": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01&$skiptoken=synthetic"
    ],
    "DisplayRegion": false
  },
  "sql-db-surrogate-name": {
    "Entries": {
      "�": {
        "State": 0,
        "BlockedZones": null
      }
    },
    "Failed": false,
    "URLs": [
      "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
    ],
    "DisplayRegion": false
  }
}`

func otherSourceOracles(t *testing.T) map[string]map[string]any {
	t.Helper()
	result := map[string]map[string]any{}
	if json.Unmarshal([]byte(otherSourceIndependentOracles), &result) != nil || len(result) != 178 {
		t.Fatal("independent remaining SKU oracle matrix invalid")
	}
	return result
}

func otherSourceRecord(t *testing.T, value any, entries any, failed bool, urls []string, display bool) {
	t.Helper()
	requests := []any{}
	for _, url := range urls {
		pathQuery := strings.TrimPrefix(url, "https://management.azure.com")
		if strings.HasPrefix(url, "https://synthetic.invalid") {
			pathQuery = strings.TrimPrefix(url, "https://synthetic.invalid")
		}
		path, query, _ := strings.Cut(pathQuery, "?")
		rawPath := ""
		if display {
			path = strings.ReplaceAll(path, "%20", " ")
			rawPath = path // Original unescaped display-region URI retained by net/http.
		}
		requests = append(requests, map[string]any{
			"URL": url, "Path": path, "RawPath": rawPath, "RawQuery": query,
			"Method": "GET", "EmptyBody": true, "SyntheticAuthorization": true,
		})
	}
	want := map[string]any{"Entries": entries, "Failed": failed, "Requests": requests}
	if !reflect.DeepEqual(value, want) {
		t.Fatal("complete remaining SKU source record mismatch")
	}
}

func TestOtherSKUSourceMetadata(t *testing.T) {
	input := computeSourceJSON(t, "source-other-sku-inputs.json")
	blobs := `{"provider.go":"7f71b296d332001ae12e49a021b6fee44f88a970","cache.go":"5be9c158fd38c609b06ffc67537cd07b838fc9d5","restrictions.go":"8d053daf58c7c8a18f39b579e95d4a53f8ebd55a","storage.go":"b5a3a8d020c8143c64d296d48e8c51ae3b0ca2cf","cognitiveservices.go":"7b3dc63af21f12213004c6e91d9eaaa06c8a2189","sql.go":"1cb1e6430f11476e66e8bbb13bc0cc2ad5e0a43f","../types/types.go":"f9b61fb74075e5a2297be1b898067b8cb047f9b3","../../../../az/http_client.go":"05d85c5a50ee81f088731160020a914af83e1ca5","compute_sku_capture_test.go":"024dd96e96c48e12a7c76831ab5a58c9e56185f5"}`
	if len(input) != 5 || !reflect.DeepEqual(input["Blobs"], computeSourceLiteral(t, blobs)) {
		t.Fatal("remaining SKU source authority mismatch")
	}
	cases, ok := input["Cases"].(map[string]any)
	if !ok || len(cases) != 178 {
		t.Fatal("remaining SKU input matrix mismatch")
	}
	for name := range otherSourceOracles(t) {
		if cases[name] == nil {
			t.Fatal("missing remaining SKU input", name)
		}
	}
	raw, err := os.ReadFile("testdata/source-other-sku-inputs.json")
	if err != nil || fmt.Sprintf("%x", sha256.Sum256(raw)) != "81ffb15038ec13e007051e79c43076a474b8a671ab4e6808cddb8d24adba3a74" {
		t.Fatal("original complete remaining SKU input bytes changed")
	}
}

func TestOtherSKUSourceCells(t *testing.T) {
	output := computeSourceJSON(t, "source-other-sku-outputs.json")
	cases, ok := output["Cases"].(map[string]any)
	if !ok || len(cases) != 178 || len(output) != 2 {
		t.Fatal("remaining SKU output matrix incomplete")
	}
	for name, want := range otherSourceOracles(t) {
		record, ok := cases[name].(map[string]any)
		if !ok || len(record) != 3 || !reflect.DeepEqual(record["Entries"], want["Entries"]) || record["Failed"] != want["Failed"] {
			t.Fatal("independent full remaining SKU map/state mismatch", name)
		}
	}
}

func TestOtherSKUSourceRequests(t *testing.T) {
	output := computeSourceJSON(t, "source-other-sku-outputs.json")
	cases, ok := output["Cases"].(map[string]any)
	if !ok || len(cases) != 178 {
		t.Fatal("remaining SKU request matrix missing")
	}
	for name, want := range otherSourceOracles(t) {
		values := want["URLs"].([]any)
		urls := make([]string, 0, len(values))
		for _, value := range values {
			urls = append(urls, value.(string))
		}
		t.Run(name, func(t *testing.T) {
			otherSourceRecord(t, cases[name], want["Entries"], want["Failed"].(bool), urls, want["DisplayRegion"].(bool))
		})
	}
}

func TestOtherSKUSourceCache(t *testing.T) {
	output := computeSourceJSON(t, "source-other-sku-outputs.json")
	cache, ok := output["Cache"].(map[string]any)
	if !ok || len(cache) != 7 {
		t.Fatal("remaining SKU cache matrix missing")
	}
	sql := "https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Sql/locations/eastus/capabilities?api-version=2021-11-01"
	urls := []string{
		sql, sql,
		strings.Replace(sql, "11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222", 1),
		strings.Replace(sql, "/eastus/", "/westus/", 1),
		"https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01",
		"https://management.azure.com/subscriptions/11111111-1111-1111-1111-111111111111/providers/Microsoft.Storage/skus?api-version=2023-01-01",
	}
	db := computeSourceLiteral(t, `{"sql_d":{"State":2,"BlockedZones":null}}`)
	mi := computeSourceLiteral(t, `{"gp_gen5":{"State":0,"BlockedZones":null}}`)
	otherSourceRecord(t, cache["db-first"], db, false, urls[:1], false)
	otherSourceRecord(t, cache["db-normalized-reuse"], db, false, urls[:1], false)
	otherSourceRecord(t, cache["mi-separate-key"], mi, false, urls[:2], false)
	otherSourceRecord(t, cache["db-subscription-separate"], db, false, urls[:3], false)
	otherSourceRecord(t, cache["db-region-separate"], db, false, urls[:4], false)
	otherSourceRecord(t, cache["cognitive-name-trim"], computeSourceLiteral(t, `{"s0":{"State":0,"BlockedZones":null}}`), false, urls[:5], false)
	otherSourceRecord(t, cache["storage-name-trim"], computeSourceLiteral(t, `{"standard_lrs":{"State":0,"BlockedZones":null}}`), false, urls, false)
}
