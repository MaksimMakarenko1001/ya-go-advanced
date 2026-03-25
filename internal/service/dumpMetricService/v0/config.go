package v0

type Config struct {
	ReadDumpEnable  bool `env:"READ_DUMP_ENABLE" envDefault:"true" json:"readDumpEnable"`
	WriteDumpEnable bool `env:"WRITE_DUMP_ENABLE" envDefault:"true" json:"writeDumpEnable"`
}
