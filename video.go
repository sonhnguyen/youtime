package youtime

func GetVideoByLink(site, id string, mongodb Mongodb) (Video, error) {
	url := URL{Site: site, ID: id}
	result, err := GetVideoByLinkMongo(url, mongodb)
	if err != nil {
		url := URL{Site: site, ID: id}
		comment := []Comment{}
		result, err = CreateVideoMongo(Video{Url: url, Comment: comment}, mongodb)
		if err != nil {
			return Video{}, err
		}
	}
	return result, nil
}
func GetVideoById(id string, mongodb Mongodb) (Video, error) {
	result, err := GetVideoByIdMongo(id, mongodb)
	if err != nil {
		return Video{}, err
	}
	return result, nil
}
func PostCommentById(id string, comment Comment, mongodb Mongodb) error {
	err := InsertCommentVideoMongo(id, comment, mongodb)
	if err != nil {
		return err
	}
	return nil
}
