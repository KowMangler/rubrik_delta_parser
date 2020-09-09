# rubrik delta parser

Simple parser for rubrik fileset delta output. Useful in scoping down to locations with heavy deltas.

Example:

`parse_diff_delta.exe --inputfile diffdata.txt --searchpath "/etc"`

```
section 1 /etc/rubrik/jobs
section 2 /etc/rubrik/upgrade_in_progress
{
   "Path": "/etc/",
   "AbsoluteMB": 10.02,
   "ReducedMB": 0.02,
   "IncreasedMB": 10,
   "TotalSizeMB": 9.98
}
```