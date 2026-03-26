package profile

import (
	"errors"

	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/internal/dto/request"
	"github.com/yourusername/autoreach-backend/internal/dto/response"
	"gorm.io/gorm"
)

func GetProfileByUserID(userID string) (*response.ProfileResponse, error) {
	var p Profile
	if err := config.DB.Where("user_id = ?", userID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toProfileResponse(p), nil
}

func CreateProfile(userID string, req request.UpsertProfileRequest) (*response.ProfileResponse, error) {
	var p Profile
	err := config.DB.Where("user_id = ?", userID).First(&p).Error
	if err == nil {
		return nil, errors.New("profile already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	p = Profile{
		UserID:      userID,
		FullName:    req.FullName,
		Title:       req.Title,
		Bio:         req.Bio,
		Location:    req.Location,
		SocialLinks: SocialLinks(req.SocialLinks),
		Skills:      req.Skills,
		ResumeRaw:   req.ResumeRaw,
		Meta:        req.Meta,
	}

	for _, e := range req.Experience {
		p.Experience = append(p.Experience, Experience(e))
	}
	for _, pr := range req.Projects {
		p.Projects = append(p.Projects, Project(pr))
	}

	if err := config.DB.Create(&p).Error; err != nil {
		return nil, err
	}

	return toProfileResponse(p), nil
}

func UpdateProfile(userID string, req request.UpsertProfileRequest) (*response.ProfileResponse, error) {
	var p Profile
	if err := config.DB.Where("user_id = ?", userID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	var experience []Experience
	for _, e := range req.Experience {
		experience = append(experience, Experience(e))
	}

	var projects []Project
	for _, pr := range req.Projects {
		projects = append(projects, Project(pr))
	}

	p.FullName = req.FullName
	p.Title = req.Title
	p.Bio = req.Bio
	p.Location = req.Location
	p.SocialLinks = SocialLinks(req.SocialLinks)
	p.Skills = req.Skills
	p.Experience = experience
	p.Projects = projects
	p.ResumeRaw = req.ResumeRaw
	p.Meta = req.Meta

	if err := config.DB.Save(&p).Error; err != nil {
		return nil, err
	}

	return toProfileResponse(p), nil
}

func toProfileResponse(p Profile) *response.ProfileResponse {
	var exp []request.ExperienceDTO
	for _, e := range p.Experience {
		exp = append(exp, request.ExperienceDTO(e))
	}

	var proj []request.ProjectDTO
	for _, pr := range p.Projects {
		proj = append(proj, request.ProjectDTO(pr))
	}

	return &response.ProfileResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		FullName:    p.FullName,
		Title:       p.Title,
		Bio:         p.Bio,
		Location:    p.Location,
		SocialLinks: request.SocialLinksDTO(p.SocialLinks),
		Skills:      p.Skills,
		Experience:  exp,
		Projects:    proj,
		ResumeRaw:   p.ResumeRaw,
		Meta:        p.Meta,
	}
}
