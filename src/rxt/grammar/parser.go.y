
{%
package grammar

import (
  "errors"

  "github.com/simelo/rextporter/src/config"
  "github.com/simelo/rextporter/src/core"
  "github.com/simelo/rextporter/src/util"
)

const (
  ErrBlockLevelUnderflow = errors.New("End of block is not possible beyond DATABASE")
  ErrBlockLevelOverflow  = errors.New("Too many nested syntax levels")
)

type parserEnv struct {
  env       core.RextEnv
  scraper   core.RextServiceScraper 
}

type strTuple struct {
  key  string
  val  string
}

type mainSecTuple struct {
  src  core.RextDataSource
  key  string
  val  interface{}
}

type metricDef {
  mname string
  mtype string
  mdesc string
  mlbls []string
  opts  core.RextKeyValueStore
}

// FIXME : Not global. Parser stack ? TLS ?
var root    parserEnv
var metric  metricDef

// TODO: metricDef should implement core.RextMetricDef

func value_for_str(str string) {
  // FIXME: Support string literals
  return string[1: len(str) - 1]
}

func newOption() core.RextKeyValueStore {
  return config.NewOptionsMap()
}

func newStrTuple(s1, s2 string) *strTuple {
  return &strTuple {
    first:  s1,
    second: s2,
  }
}

func newMainDef(key string, value inrerface{}) *mainSecTuple {
  return &mainSecTuple{
    src: nil,
    key: key,
    value: value,
  }
}

func newMainSrc(src core.RextDataSource) *mainSecTuple {
  return &mainSecTuple{
    src: src,
    key: "",
    value: nil,
  }
}

func getRootEnv() *parserEnv {
  return &root
}

func (m *metricDef) GetMetricName() string {
  return m.mname
}

func (m *metricDef) GetMetricType() string {
  return m.mtype
}

func (m *metricDef) GetMetricDescription() string {
  return m.mdesc
}

func (m *metricDef) GetMetricLabels() []string {
  return m.mlbls
}

func (m *metricDef) SetMetricName(name string) {
  m.mname = name
}

func (m *metricDef) SetMetricType(typeid string) {
  m.mtype = typeid
}

func (m *metricDef) SetMetricDescription(desc string) {
  m.mdesc = desc
}

func (m *metricDef) SetMetricLabels(labels []string) {
  m.mlbls = labels
}

func (m *metricDef) GetOptions() RextKeyValueStore {
  return nil
}

%}

%union{
  root    core.RextServiceScraper
  options core.RextKeyValueStore
  mains   []mainSecTuple
  mainsec *mainSecTuple
  exts    []core.RextMetricsExtractor
  extract core.RextMetricsExtractor
  metrics []metricsDef
  metric  metricsDef
  key     string
  strval  string
  strlist []string
  pair    *strTuple
}

%type <strval>  id mname mtype mhelp mhelpo
%type <key>     defverb srcverb mtvalue mfname
%type <pair>    setcls
%type <strlist> strlst idlst mlabels mlablso stkcls stkclso srvcls srvclso
%type <options> optsblk optblkl optblkr optblko
%type <metric>  metsec
%type <metrics> metblk
%type <extract> extblk
%type <mainsec> srcsec defsec mainsec
%type <mains>   mainblk
%type <root>    dataset

%token <strval> STR VAR

%%

defverb : 'DEFINE AUTH'
          { $$ = "AUTH" }
        ;
srcverb : 'GET'
          { $$ = $1 }
        | 'POST'
          { $$ = $1 }
        ;
mtvalue : 'GAUGE'
          { $$ = config.KeyTypeGauge }
        | 'COUNTER'
          { $$ = config.KeyTypeCounter }
        | 'HISTOGRAM'
          { $$ = config.KeyTypeHistogram }
        | 'SUMMARY'
          { $$ = config.KeyTypeSummary }
        ;
id      : VAR
          { $$ = $1 }
        | STR
          { $$ = value_for_str($1) }
        ;
setcls  : 'SET' STR 'TO' STR
          { $$ = newStrTuple($2, $4) }
optsblk : setcls
          {
            // TODO: Error handling
            $$ = newOption()
            _, _ = $$.SetString($1.first, $1.second)
          }
        | optsblk EOL setcls
          {
            // TODO: Error handling
            _, _ = $1.top().SetString($3.first, $3.second)
            $$ = $1
          }
strlst : STR
          { $$ = []string{ $1 } }
        | strlst ',' STR
          { $$ = append($1, $3) }
idlst   : id
          { $$ = []string{ $1 } }
        | idlst ',' id
          { $$ = append($1, $3) }
mlabels : 'LABELS' strlst
          { $$ = $2 }
mname   : 'NAME' ID
          { $$ = $2 }
mtype   : 'TYPE' mtvalue
          { $$ = $2 }
mhelp   : 'HELP' STR
          { $$ = $2 }
mhelpo  : /* empty */
          {
            $$ = "Metric extracted by [rextporter](https://github.com/simelo/rextporter)"
          }
        | EOL mhelp
          { $$ = $2 }
mlablso : /* empty */
          { $$ = nil }
        | EOL mlabels
          { $$ = $2 }
optblkl : /* empty */
          { $$ = nil }
          | EOL optsblk
          { $$ = $2 }
optblkr : /* empty */
          { $$ = nil }
          | optsblk EOL
          { $$ = $1 }
metsec  : 'METRIC' BLK mname EOL mtype mhelpo mlablso optblkl EOB
          {
            $$ = metricDef{
              mname: $3,
              mtype: $5,
              mdesc: $6,
              mlbls: $7,
              opts:  $8,
            }
          }
metblk  : metsec
          { $$ = []metricsDef{ $1 } }
        | metblk EOL metsec
          { $$ = append($1, $3) }
extblk  : 'EXTRACT USING' id BLK optblkr metblk EOB
          {
            env := getRootEnv()
            $$ = env.NewMetricsExtractor($2, $4, $5)
            for _, md := range $6 {
              $$.AddMetricRule(&md)
            }
          }
ssec    : extblk
          { $$ = []core.RextMetricsExtractor{ $1 } }
        | ssec EOL extblk
          { $$ = append($1, $2) }
srcsec  : srcverb VAR 'FROM' STR
          {
            env := getRootEnv()
            ds := env.NewMetricsDataSource($2)
            dsSetMethod($1)
            dsSetLocation($4)
            $$ = newMainSrc(ds)
          }
        | srcverb VAR 'FROM' STR BLK optblkr ssec EOB
          {
            env := getRootEnv()
            ds := env.NewMetricsDataSource($2)
            ds.SetMethod($1)
            ds.SetLocation($4)
            // FIXME: Error handling
            _ = util.MergeStoresInplace(dsGetOptions(), $6)
            $$ = newMainSrc(ds)
          }
defsec  : defverb VAR 'AS' id optblko
          {
            env := getRootEnv()
            if defverb == 'AUTH' {
              $$ = newMainDef($4, env.NewAuthStrategy($2, $5))
            }
            // TODO: Error handling
            $$ = nil
          }
optblko : /* empty */
          { $$ = nil }
        | BLK optsblk EOB
          { $$ = $2 }
stkcls  : 'FOR STACK' idlst
          { $$ = $2 }
stkclso : /* empty */
          { $$ = nil }
        | stkcls EOL
          { $$ = $1 }
srvcls  : 'FOR SERVICE' idlst
          { $$ = $2 }
srvclso : /* empty */
          { $$ = nil }
        | srvcls EOL
          { $$ = $1 }
mainsec : defsec
          { $$ = $1 }
        | srcsec
          { $$ = $1 }
mainblk : mainsec
          { $$ = []mainSecTuple { $1 } }
        | mainblk EOL mainsec
          { $$ = append($1, $2) }
eolo    : /* empty */
        | EOL
dataset : CTX eolo 'DATASET' BLK srvclso srvstko optblkr mainblk EOB eolo
          {
            env = $1
            $$ = env.NewServiceScraper()
            if $5 != nil {
              // TODO : Error handling
              _ = env.RegisterScraperForServices($5...)
            }
            if $6 != nil {
              // TODO : Error handling
              _ = env.RegisterScraperForServices($6...)
            }
            if $7 != nil {
              util.MergeStoresInplace($$.GetOptions(), $7)
            }
            for _, mainsec := range $8 {
              if mainsec.src != nil {
                $$.AddSource(mainsec.src)
              } else if mainsec.value != nil {
                if auth, isAuth := mainsec.value.(core.RextAuth); isAuth {
                  $$.AddAuthStrategy(auth, mainsec.key)
                }
                // TODO : Error handling
              }
              // TODO : Error handling
            }
          }
%%






