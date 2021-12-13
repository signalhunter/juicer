# Juicer

It's ripgrep but for Gzip-compressed files over HTTP!

This tool was primarily designed to scan thru the Common Crawl dataset for URLs without spending a fortune on AWS.

Features:
  - Extremely fast regex engine ([Intel Hyperscan](https://www.hyperscan.io/))
  - Scan thru terabytes of data without writing them to disk
  - Concurrent scanning of multiple files

TODO:
  - Client/server for handing out scanning tasks
  - Zstandard support? (for IA WARCs)
