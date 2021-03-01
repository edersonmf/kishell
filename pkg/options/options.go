package options

import (
  "github.com/alecthomas/kong"
  "github.com/sidilabs/kishell/pkg/config"
  "github.com/sidilabs/kishell/pkg/utils"
  "time"
)

type Context struct {
  Debug         bool
  Configuration config.Configuration
}

type Option struct {
  Context *kong.Context
  ConfigFile config.Configuration
}

type ConfigureCmd struct {
  Server bool `optional help:"Add a new server definition"`
  Role bool `optional help:"Add a new role definition"`
  Reset bool `optional help:"Reset the whole configuration"`
}

type UseCmd struct {
  Server string `optional help:"Set what server to use when querying ES"`
  Role string `optional help:"Set what role to use when querying ES"`
}

type ListCmd struct {
}

type SearchCmd struct {
  Query string `optional help:"Text input to query data. Use the same format as you would use in Kibana"`
  Older string `optional default:"now" help:"Data older than. Defaults to current time when not provided (e.g. 30m, 1h, 1w, 1M, 1y)"`
  Newer string `optional default:"15m" help:"Data newer than. Defaults to 15m when not provided (e.g. 30m, 1h, 1w, 1M, 1y)"`
  Limit int32 `optional default:"50" help:"Limit the number of messages fetched"`
  Server string `optional help:"Which server to query against. Used to override the current server config"`
  httpClient utils.HttpClient `-`
}

var CLI struct {
  Debug bool `help:"Enable debug mode."`
  Configure ConfigureCmd `cmd help:"Init ES server configs"`
  List ListCmd `cmd help:"Show the current server configs"`
  Search SearchCmd `cmd help:"Search for data"`
  Use UseCmd `cmd help:"Update config options with ser/role preferences"`
}

func (s *SearchCmd) OlderAsTimestamp() (int64, error) {
  return toTimestamp(s.Older)
}

func (s *SearchCmd) NewerAsTimestamp() (int64, error) {
  return toTimestamp(s.Newer)
}

func (s *SearchCmd) AfterApply(h *utils.DefaultHttpClient) error {
  s.httpClient = h
  return nil
}

func toTimestamp(period string) (int64, error) {
  now := time.Now().Unix() * 1000
  if len(period) <= 0 || period == "now" {
    return now, nil
  }
  duration, err := time.ParseDuration(period)
  if err != nil {
    return -1, err
  }
  return now - duration.Milliseconds(), nil
}

func (o *Option) Run() {
  err := o.Context.Run(&Context {
      Debug:         CLI.Debug,
      Configuration: o.ConfigFile,
  })
  o.Context.FatalIfErrorf(err)
}

func Parse() Option {
  httpClient := &utils.DefaultHttpClient{
    Timeout: 5 * time.Second,
  }
  context := kong.Parse(&CLI, kong.Bind(httpClient))
  opt := Option {
    Context: context,
    ConfigFile: config.LoadDefaultConfig(),
  }
  return opt
}
