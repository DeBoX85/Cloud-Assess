# Additive history consumer coverage

Static review after PR211 own QA found three gaps in independent coverage of already frozen requirements. No producer exists and no mutation defect is claimed. Preserve the eight original files and all 77 original future events unchanged.

Before any producer code, publish/read back these separate files, observe actual missing APIs and compile the fixture-only shape without invoking an API. Add three roots/five named events: uppercase USD request versus returned usd must reject (RH_CURRENCY_EXACT); 256 copies of é are exactly 512 UTF8 bytes and must be accepted, 257 copies are 514 bytes and reject (RH_UTF8_BYTES); empty name on selected non-requested subscription B must reject even when request A is valid (RH_SCOPE_ALL_NAMES). Total future events are 82 in 17 roots. No registry refresh, currency conversion or runtime mutation credit follows from compiler failure.

Copy consumers unchanged into implementation. New strict reversals must compile, fail at the intended marker, restore exactly and pass, in addition to all original controls. Own full native/source QA and private retention remain required. No coordinator/task10 closure.
