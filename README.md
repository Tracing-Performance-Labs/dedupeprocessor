# Dedupe Procesor

| Status | |
|------------|---|
| Stability | Alpha |
| Code Owners | @aloussase | 

The dedupe processor deduplicates string data in span attributes. It should be
used as the last processor in the pipeline if any other processing on
attributes is required.

Applications that use this processor should install a decoding proxy in front
of their backends so that the original data can be recovered.
