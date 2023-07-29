package host

import "fmt"

type HostURL struct {
	First_Chunk  int // 127
	Second_Chunk int // 0-255
	Third_Chunk  int // 0-255
	Fourth_Chunk int // 0-255
}

func NewHostURL() *HostURL {
	return &HostURL{
		First_Chunk:  127,
		Second_Chunk: 0,
		Third_Chunk:  0,
		Fourth_Chunk: 0,
	}
}

func (h *HostURL) UpgradeHostURL() error {

	if h.Fourth_Chunk < 255 {
		h.Fourth_Chunk++
		return nil
	} else if h.Third_Chunk < 255 {
		h.Third_Chunk++
		h.Fourth_Chunk = 1
		return nil
	} else if h.Second_Chunk < 255 {
		h.Second_Chunk++
		h.Third_Chunk = 0
		h.Fourth_Chunk = 0
		return nil
	} else {
		return fmt.Errorf("HostURL is out of range")
	}

}

func (h *HostURL) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", h.First_Chunk, h.Second_Chunk, h.Third_Chunk, h.Fourth_Chunk)
}
