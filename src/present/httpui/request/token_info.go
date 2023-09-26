package request

type TokenInfo struct {
	ShopCode        string //kvr code
	IsAdmin         bool   //kvuadmin
	RetailerId      int64  //kvrid
	RetailerUserId  int64  //kvuid
	RetailerUser    string //preferred_username
	UsernameUser    string //preferred_username
	Token           string //token
	BranchId        int    //header branch-id or kvbid
	VersionLocation int    //header version-location, default = 1
}
