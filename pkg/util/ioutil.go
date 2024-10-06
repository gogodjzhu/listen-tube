package util

type ChanWriter chan string

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	str := string(p)   // 将字节切片转换为字符串
	cw <- str          // 将字符串写入通道
	return len(p), nil // 返回写入的字节数
}
