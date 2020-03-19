# United States Death Master File 2 CSV Converter

For details read this blog post: https://blog.intelx.io/2020/03/19/decoding-the-us-death-master-file/

Here's the converted file indexed in our search engine: https://intelx.io/?did=fd36a1b3-35ff-429b-8c19-e6a4e229ffb9

## Compile

Download and install Go from https://golang.org/dl/. Compile it using this command (it compiles on Windows, Linux and Mac):

```
go build
```

## Example input and output

Input:

```
 001010001MUZZEY                  GRACE                          1200197504161902                   
 001010009SMITH                   ROGER                          0400196902041892                   
 001010010HAMMOND                 KENNETH                        0300197604241904                   
 001010011DREW                    LEON           R              V0830198706141908                   
```

Output as CSV:

```
Type,Social Security Number,Last Name,Name Suffix,First Name,Middle Name,Verified,Date of Death,Date of Birth,Blank 1,Blank2,Blank 3,Blank 4
,001010001,Muzzey,,Grace,,,1975-12-00,1902-04-16,,,,
,001010009,Smith,,Roger,,,1969-04-00,1892-02-04,,,,
,001010010,Hammond,,Kenneth,,,1976-03-00,1904-04-24,,,,
,001010011,Drew,,Leon,R,Verified,1987-08-30,1908-06-14,,,,
```

## License

This is free and unencumbered software released into the public domain.
