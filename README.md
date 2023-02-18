## go-cms-dl
[![Go Report Card](https://goreportcard.com/badge/github.com/mahmednabil109/go-cms-dl)](https://goreportcard.com/report/github.com/mahmednabil109/go-cms-dl)
![Lines of Code](https://tokei.rs/b1/github/mahmednabil109/go-cms-dl?category=code)

simple web scraping go based tool to download __a__ cms courses content, with support of concurrent http connections, which is completely useless due to the fact that the cms __awesome__ server has a strict rate limiter 😀.
## TODOS
- [x] modulize the code.
- [x] make project configurable (yaml).
- [x] presist some metadata about downloaded files (optional).
- [ ] convert to [cobra](https://github.com/spf13/cobra). 
- [ ] serve simple html page displaying some stats about the downloaded files.
