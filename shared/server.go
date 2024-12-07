package shared

// サーバーの起動を共通化するためのインターフェース
type Server interface {
	// Run サーバーの起動
	Run(addr string)
}
