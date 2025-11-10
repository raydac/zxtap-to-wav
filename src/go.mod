module github.com/raydac/zxtap-to-wav

go 1.25

require (
        github.com/raydac/zxtap-wav v1.0.1
        github.com/raydac/zxtap-zxtape v1.0.1
        github.com/raydac/zxtap-zx v1.0.1
)

replace (
        github.com/raydac/zxtap-wav => ./wav
        github.com/raydac/zxtap-zxtape => ./zxtape
        github.com/raydac/zxtap-zx => ./zx
)