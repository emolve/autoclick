package setting

type NotificationSettingS struct {
	Mail     string
	PlusPlus string
}
type AppSettingS struct {
	RunMode string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
