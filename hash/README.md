# Install

```
make build
make integration-test
```

# Use

```
# bin/hash
NAME:
   hash - hash command [OPTIONS] filename

USAGE:
   hash [global options] command [command options] [arguments...]

VERSION:
   0.9.9

COMMANDS:
     md5      hash md5 [OPTIONS] filename
     sha1     hash sha1 [OPTIONS] filename
     sha256   hash sha256 [OPTIONS] filename
     sha512   hash sha512 [OPTIONS] filename
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the versio

```

```
# bin/hash md5
invalid filename

NAME:
   hash md5 - hash md5 [OPTIONS] filename

USAGE:
   hash md5 [command options] [arguments...]

OPTIONS:
   -R                         hash entire subtree connected at that point.
   -U, --toUpper              output upper(default is lower)
   -H, --human                Human-readable
   -o value, --output value   output file
   -f value, --format value   output format, use golang template (Name, FullName, Size, Hash)
   --hmac value               hmachash
   -E value, --exclude value  exclude

```

eq.
```
# bin/hash md5 -R -H --exclude "***_test.go" --exclude "vendor" --exclude ccc.goo ./
./.DS_Store, 6.0K, 47bc6d4786532deaec6e59bb9f8d0930
./Makefile, 606B, c003f78f24eb35759a2902db7806dc47
./README.md, 1.1K, 33c2d1f3c743a97f7e2b172a9380b154
./VERSION.txt, 6B, e81aaac78a68b93fc4ae622c88fac606
./bin/hash, 3.6M, 57d73720f8da6296663ea8fe0d7a17f3
./hash.go, 2.7K, f8915c7099f7bd052f11233e79814fe6
./main.go, 2.4K, ae1a474df958538a5e80019d9b7500ec
./tests/a.txt, 6B, 4c850c5b3b2756e67a91bad8e046ddac
./tests/check.sh, 104B, cf7442cd2df084390afbb0a34926932e
./tests/check_darwin.go, 596B, 308e24f7eb0577dfac090485cd59ed2f
./tests/check_linux.go, 569B, 140151c351a2c36506c9e1ac45eabef1
./types.go, 1.3K, 7dff4660e35840854ec4b0fcbd6c5779
./util.go, 648B, 134545050d76d8bf9f52a770aa08fae8
```
