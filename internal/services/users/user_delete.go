package users

func (s *Service) Delete(id int) error {
	return s.Repo.DeleteUser(id)
}
