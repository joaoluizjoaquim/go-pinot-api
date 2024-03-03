package model

type TableSegmentsConfig struct {
	TimeType                  string `json:"timeType"`
	Replication               string `json:"replication"`
	TimeColumnName            string `json:"timeColumnName"`
	SegmentAssignmentStrategy string `json:"segmentAssignmentStrategy"`
	SegmentPushType           string `json:"segmentPushType"`
	MinimizeDataMovement      bool   `json:"minimizeDataMovement"`
}

type TableTenant struct {
	Broker string `json:"broker"`
	Server string `json:"server"`
}

type StarTreeIndexConfig struct {
	DimensionsSplitOrder              []string `json:"dimensionsSplitOrder"`
	SkipStarNodeCreationForDimensions []string `json:"skipStarNodeCreationForDimensions"`
	FunctionColumnPairs               []string `json:"functionColumnPairs"`
	MaxLeafRecords                    int      `json:"maxLeafRecords"`
}

type TierOverwrite struct {
	StarTreeIndexConfigs []StarTreeIndexConfig `json:"starTreeIndexConfigs"`
}

type TierOverwrites struct {
	HotTier  TierOverwrite `json:"hotTier"`
	ColdTier TierOverwrite `json:"coldTier"`
}

type TableIndexConfig struct {
	EnableDefaultStarTree                      bool                  `json:"enableDefaultStarTree"`
	StarTreeIndexConfigs                       []StarTreeIndexConfig `json:"starTreeIndexConfigs"`
	TierOverwrites                             TierOverwrites        `json:"tierOverwrites"`
	EnableDynamicStarTreeCreation              bool                  `json:"enableDynamicStarTreeCreation"`
	AggregateMetrics                           bool                  `json:"aggregateMetrics"`
	NullHandlingEnabled                        bool                  `json:"nullHandlingEnabled"`
	OptimizeDictionary                         bool                  `json:"optimizeDictionary"`
	OptimizeDictionaryForMetrics               bool                  `json:"optimizeDictionaryForMetrics"`
	NoDictionarySizeRatioThreshold             float64               `json:"noDictionarySizeRatioThreshold"`
	RangeIndexVersion                          int                   `json:"rangeIndexVersion"`
	AutoGeneratedInvertedIndex                 bool                  `json:"autoGeneratedInvertedIndex"`
	CreateInvertedIndexDuringSegmentGeneration bool                  `json:"createInvertedIndexDuringSegmentGeneration"`
	LoadMode                                   string                `json:"loadMode"`
}

type TableMetadata struct {
	CustomConfigs map[string]string `json:"customConfigs"`
}

type TimestampConfig struct {
	Granulatities []string `json:"granularities"`
}

type FiendIndexInverted struct {
	Enabled string `json:"enabled"`
}

type FieldIndexes struct {
	Inverted FiendIndexInverted `json:"inverted"`
}

type FieldConfig struct {
	Name            string          `json:"name"`
	EncodingType    string          `json:"encodingType"`
	IndexType       string          `json:"indexType"`
	IndexTypes      []string        `json:"indexTypes"`
	TimestampConfig TimestampConfig `json:"timestampConfig"`
	Indexes         FieldIndexes    `json:"indexes"`
}

type TransformConfig struct {
	ColumnName        string `json:"columnName"`
	TransformFunction string `json:"transformFunction"`
}

type TableIngestionConfig struct {
	SegmentTimeValueCheckType string            `json:"segmentTimeValueCheckType"`
	TransformConfigs          []TransformConfig `json:"transformConfigs"`
	ContinueOnError           bool              `json:"continueOnError"`
	RowTimeValueCheck         bool              `json:"rowTimeValueCheck"`
}

type TierConfig struct {
	Name                string `json:"name"`
	SegmentSelectorType string `json:"segmentSelectorType"`
	SegmentAge          string `json:"segmentAge"`
	StorageType         string `json:"storageType"`
	ServerTag           string `json:"serverTag"`
}

type Table struct {
	TableName        string               `json:"tableName"`
	TableType        string               `json:"tableType"`
	SegmentsConfig   TableSegmentsConfig  `json:"segmentsConfig"`
	Tenants          TableTenant          `json:"tenants"`
	TableIndexConfig TableIndexConfig     `json:"tableIndexConfig"`
	Metadata         TableMetadata        `json:"metadata"`
	FieldConfigList  []FieldConfig        `json:"fieldConfigList"`
	IngestionConfig  TableIngestionConfig `json:"ingestionConfig"`
	TierConfigs      []TierConfig         `json:"tierConfigs"`
	IsDimTable       bool                 `json:"isDimTable"`
}
