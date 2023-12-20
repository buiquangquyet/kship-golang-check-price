package enums

type TypeService int

const (
	TypeServiceDV   TypeService = 1
	TypeServiceDVMR TypeService = 2
)

func (m TypeService) ToInt() int {
	return int(m)
}
