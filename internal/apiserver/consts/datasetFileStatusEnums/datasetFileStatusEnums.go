package datasetfilestatusenums

var (
	UnLabel  uint32 = 0   //待标注
	Labeling uint32 = 10  //标注中
	Labeled  uint32 = 100 //已标注
	Deleted  uint32 = 255 //标记删除
)
