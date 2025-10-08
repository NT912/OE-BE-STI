package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type AiService struct{}

func NewAiService() *AiService {
	return &AiService{}
}

func (s *AiService) GetAIStream(writer io.Writer) error {
	responseWords := []string{
		"Trong ", "một ", "thế ", "giới ", "ngày ", "càng ", "phát ", "triển, ", "trí ", "tuệ ", "nhân ", "tạo ",
		"đã ", "trở ", "thành ", "một ", "phần ", "thiết ", "yếu ", "trong ", "cuộc ", "sống ", "của ", "chúng ", "ta. ",
		"Từ ", "những ", "ứng ", "dụng ", "đơn ", "giản ", "như ", "gợi ", "ý ", "từ ", "khóa, ",
		"cho ", "đến ", "các ", "hệ ", "thống ", "phức ", "tạp ", "có ", "khả ", "năng ", "học ", "hỏi ", "và ", "thích ", "nghi, ",
		"AI ", "đang ", "dần ", "thay ", "đổi ", "cách ", "con ", "người ", "làm ", "việc ", "và ", "sáng ", "tạo. ",
		"Khi ", "bạn ", "đọc ", "những ", "dòng ", "chữ ", "này, ", "rất ", "có ", "thể ", "một ", "mô ", "hình ", "AI ",
		"đang ", "xử ", "lý ", "hàng ", "triệu ", "tác ", "vụ ", "trong ", "một ", "giây, ", "học ", "từ ", "hàng ", "tỷ ", "dữ ", "liệu ",
		"để ", "hiểu ", "về ", "ngôn ", "ngữ, ", "về ", "con ", "người, ", "và ", "về ", "thế ", "giới ", "xung ", "quanh. ",
		"Điều ", "thú ", "vị ", "là, ", "những ", "phản ", "hồi ", "bạn ", "nhận ", "được ", "ở ", "đây ",
		"không ", "phải ", "chỉ ", "là ", "chuỗi ", "ký ", "tự, ", "mà ", "là ", "kết ", "quả ", "của ",
		"hàng ", "triệu ", "tính ", "toán ", "song ", "song ", "được ", "thực ", "hiện ", "chỉ ", "trong ", "chớp ", "mắt. ",
		"Và ", "tất ", "cả ", "những ", "điều ", "đó ", "đang ", "được ", "stream ", "từng ", "chút ",
		"một, ", "đến ", "màn ", "hình ", "của ", "bạn, ", "theo ", "thời ", "gian ", "thực. ",
		"Đây ", "là ", "sức ", "mạnh ", "của ", "AI ", "thật ", "sự.",
	}

	flusher, ok := writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	for _, word := range responseWords {
		_, err := fmt.Fprint(writer, word)
		if err != nil {
			return err
		}
		flusher.Flush()
		log.Printf("Sent Chunk: %s", word)
		// time.Sleep(200 * time.Millisecond)
	}
	log.Println("Finished streaming AI response.")
	return nil
}
