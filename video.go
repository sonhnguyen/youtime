package youtime

func GetVideoByLink(site, link string, mongodb Mongodb) (Video, error) {
	url := URL{Site: site, Link: link}
	result, err := GetVideoByLinkMongo(url, mongodb)
	if err != nil {
		return Video{}, err
	}
	if result.Id == "" {
		url := URL{Site: site, Link: link}
		commend := []Comment{}
		result, err = CreateVideoMongo(Video{Url: url, Comment: commend}, mongodb)
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
