package handlers

// type inputCategory struct {
// 	Name      string `json:"name"`
// 	Priority  uint   `json:"priority"`
// 	Recurrent bool   `json:"recurrent"`
// 	Notify    bool   `json:"notify"`
// }

// func (s *Store) CreateCategory(w http.ResponseWriter, r *http.Request) {
// 	inputCate := new(inputCategory)
// 	category := new(models.Category)

// 	defer r.Body.Close()
// 	if err := json.NewDecoder(r.Body).Decode(&inputCate); err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	if inputCate.Name == "" {
// 		http.Error(w, "a name is required", http.StatusBadRequest)
// 		return
// 	}

// 	message, userID := getUserId(w, r)
// 	if userID == -1 {
// 		http.Error(w, message, http.StatusBadRequest)
// 		return
// 	}

// 	category.UserID = uint(userID)
// 	category.Name = inputCate.Name
// 	category.Priority = inputCate.Priority
// 	category.Recurrent = inputCate.Recurrent
// 	category.Notify = inputCate.Notify

// 	if err := s.DB.Create(&category).Error; err != nil {
// 		http.Error(w, "error creating category", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(category)
// }

// func (s *Store) GetCategories(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	message, userID := getUserId(w, r)

// 	if userID == -1 {
// 		http.Error(w, message, http.StatusBadRequest)
// 		return
// 	}

// 	if id != "" {
// 		s.getUserFromId(id, w)
// 	}
// 	s.getAllCategories(userID, w)
// }

// func (s *Store) getAllCategories(userID int, w http.ResponseWriter) {
// 	categories := []models.Category{}

// 	if err := s.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
// 		http.Error(w, "error getting categories", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	if err := json.NewEncoder(w).Encode(categories); err != nil {
// 		http.Error(w, "error encoding categories", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (s *Store) getCategoryFromId(id string, w http.ResponseWriter) {
// 	var category models.Category

// 	if err := s.DB.First(&category, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "category not found", http.StatusNotFound)
// 			return
// 		}

// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(category)
// }

// func (s *Store) UpdateCategory(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	var updatedCategory inputCategory
// 	var category models.Category

// 	defer r.Body.Close()
// 	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	if err := s.DB.First(&category, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "category not found", http.StatusNotFound)
// 			return
// 		}
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	if updatedCategory.Name == "" {
// 		http.Error(w, "empty fields", http.StatusBadRequest)
// 		return
// 	}

// 	category.Name = updatedCategory.Name
// 	category.Priority = updatedCategory.Priority
// 	category.Recurrent = updatedCategory.Recurrent
// 	category.Notify = updatedCategory.Notify

// 	if err := s.DB.Save(&category).Error; err != nil {
// 		if errors.Is(err, gorm.ErrInvalidData) {
// 			http.Error(w, "invalid data", http.StatusBadRequest)
// 			return
// 		} else {
// 			http.Error(w, "internal error", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(category)
// }

// func (s *Store) DeleteCategory(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	var category models.Category

// 	if err := s.DB.First(&category, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			http.Error(w, "category not found", http.StatusNotFound)
// 			return
// 		}
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := s.DB.Delete(&category).Error; err != nil {
// 		http.Error(w, "error deleting", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	response := map[string]string{"message": "category deleted successfully"}
// 	json.NewEncoder(w).Encode(response)
// }
