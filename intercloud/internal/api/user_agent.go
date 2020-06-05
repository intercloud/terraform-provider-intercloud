package api

type UserAgentProduct struct {
	Name    string
	Version string
}

func (u UserAgentProduct) String() (str string) {
	str = u.Name + "/" + u.Version
	return
}

type UserAgentProducts []UserAgentProduct

func (u UserAgentProducts) String() (str string) {
	for i := range u {
		str += " " + u[i].String()
	}
	if str != "" {
		str = str[1:]
	}

	return
}
