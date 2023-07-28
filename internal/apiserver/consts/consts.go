package consts

const TimeFormat = "2006-01-02 15:04:05"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
	RoleView  = "view"
)

const (
	NotUsed = iota // 未使用
	Used           // 已使用
)

const (
	TokenRegister = iota + 1 // token 为注册类型
	TokenForget              // token 为忘记密码
)

const (
	TrainRecordStatusCreated       = 0
	TrainRecordStatusLabelManaged  = 10
	TrainRecordStatusDataUploaded  = 20
	TrainRecordStatusLabelsMarked  = 30
	TrainRecordStatusModelTraining = 40
	TrainRecordStatusTrainFailed   = 50
	TrainRecordStatusTrainSuccess  = 60
	TrainRecordStatusDeploying     = 70
	TrainRecordStatusDeployFailed  = 80
	TrainRecordStatusDeploySuccess = 90
	TrainRecordStatusModelReady    = 100
	TrainRecordStatusModelDeleted  = 255
	TrainRecordStatusJobDeleted    = 404
	TrainJobNameRegex              = "train-%d-%d-%d-%s"
	DeployEngineNameRegex          = "engine-%d-%d-%d"
	TrainJobNameRegexExtract       = `^train-(\d+)-(\d+)-(.*)$`
)

const (
	EngineStatusDeleted = 0
	EngineStatusNormal  = 1
)

const (
	EngineDeployConfigModelUrlPrefix = "http://camfs-service.camfs.svc.cluster.local/download/aws/cn-northwest-1/textin-cn-test/"
	CallBackURL                      = "http://textin-global-backend.default.svc.cluster.local:8080/internal/api/v1/model/"
)

const (
	ListTaskTicker = 30
)

const (
	ApiFree = iota
	ApiPay
)
