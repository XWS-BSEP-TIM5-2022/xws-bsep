package api

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/startup/config"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	verificationCodeDurationInMinutes int = 5
	min6DigitNumber                   int = 100000
	max6DigitNumber                   int = 999999
	minPasswordLength                 int = 8
)

var validate *validator.Validate

type AuthService struct {
	store             *persistence.AuthPostgresStore
	jwtService        *JWTService
	apiTokenService   *APITokenService
	userServiceClient user.UserServiceClient
	CustomLogger      *CustomLogger
}

type LoginData struct {
	username string `validate:"required"`
	password string `validate:"required"`
}

func NewAuthService(store *persistence.AuthPostgresStore, jwtService *JWTService, userServiceClient user.UserServiceClient, apiTokenService *APITokenService) *AuthService {
	CustomLogger := NewCustomLogger()
	return &AuthService{
		store:             store,
		jwtService:        jwtService,
		userServiceClient: userServiceClient,
		apiTokenService:   apiTokenService,
		CustomLogger:      CustomLogger,
	}
}

func (service *AuthService) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {
	service.CustomLogger.InfoLogger.Info("Passwordless login for user with email: " + request.Email)
	err := checkEmailCriteria(request.Email)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("Email: " + request.Email + " is invalid")
		fmt.Println(err.Error())
		return nil, err
	}
	getUserRequest := &user.GetIdByEmailRequest{
		Email: request.Email,
	}

	user, err := service.userServiceClient.GetIdByEmail(context.TODO(), getUserRequest)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("No user with email: " + request.Email + " or account is not activated")
		return nil, errors.New("there is no user with that email or account is not activated")
	}

	authCredentials, err := service.store.FindById(user.Id)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("No user found with ID: " + user.Id)
		return nil, errors.New("user not found")
	}

	service.CustomLogger.DebugLogger.Info("Finding roles from user with ID: " + user.Id)
	var authRoles []domain.Role
	for _, authRole := range *authCredentials.Roles {
		roles, err := service.store.FindRoleByName(authRole.Name)
		if err != nil {
			service.CustomLogger.ErrorLogger.Error("No role found with name: " + authRole.Name)
			fmt.Println("Error finding role by name")
			return nil, err
		}
		authRoles = append(authRoles, *roles...)
	}
	authCredentials.Roles = &authRoles

	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("JWT token is not generated for user with ID: " + user.Id)
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}

	service.CustomLogger.DebugLogger.Info("Sending passwordless login email for user with ID: " + user.Id)
	message, subject := passwordlessLoginMailMessage(token)

	err = sendEmail(request.Email, message, subject)
	if err != nil {
		fmt.Println(err)
		service.CustomLogger.ErrorLogger.Error("Passwordless login email not sent to user with ID: " + user.Id)
		return nil, errors.New("error while sending mail")
	}

	service.CustomLogger.SuccessLogger.Info("Passwordless login email successfully sent to user with ID: " + user.Id)
	return &pb.PasswordlessLoginResponse{
		Success: "Email sent successfully! Check your email.",
	}, nil
}

func passwordlessLoginMailMessage(token string) (string, string) {
	urlRedirection := "https://" + config.NewConfig().FrontendHost + ":" + config.NewConfig().FrontendPort + "/confirmed-mail/" + token

	subject := "Passwordless login"
	body := "<html><body style=\"background-color: #f4f4f4; margin: 0 !important; padding: 0 !important;\">\n" +
		"    <!-- HIDDEN PREHEADER TEXT -->\n" +
		"    <div style=\"display: none; font-size: 1px; color: #fefefe; line-height: 1px; font-family: 'Lato', Helvetica, Arial, sans-serif; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden;\"> We're thrilled to have you here! Get ready to dive into your new account.\n" +
		"    </div>\n" +
		"    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\">\n" +
		"        <!-- LOGO -->\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td align=\"center\" valign=\"top\" style=\"padding: 40px 10px 40px 10px;\"> </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"center\" valign=\"top\" style=\"padding: 40px 20px 20px 20px; border-radius: 4px 4px 0px 0px; color: #111111; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 48px; font-weight: 400; letter-spacing: 4px; line-height: 48px;\">\n" +
		"                            <h1 style=\"font-size: 48px; font-weight: 400; margin: 2;\">Dislinkt</h1> <img src=\" https://img.icons8.com/cotton/100/000000/security-checked--v3.png\" width=\"125\" height=\"120\" style=\"display: block; border: 0px;\" />\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#f4f4f4\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\" style=\"padding: 20px 30px 40px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 25px;\">\n" +
		"                            <p style=\"margin: 0;\">Someone tried to sign in to your account without password. Was that you?</p>\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\">\n" +
		"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                <tr>\n" +
		"                                    <td bgcolor=\"#ffffff\" align=\"center\" style=\"padding: 20px 30px 60px 30px;\">\n" +
		"                                        <table border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                            <tr>\n" +
		"                                                <td align=\"center\" style=\"border-radius: 3px;\" bgcolor=\"#FFA73B\"><a href=\"" + urlRedirection + "\" target=\"_blank\" style=\"font-size: 20px; font-family: Helvetica, Arial, sans-serif; color: #ffffff; text-decoration: none; color: #ffffff; text-decoration: none; padding: 15px 25px; border-radius: 2px; border: 1px solid #FFA73B; display: inline-block;\">Yes! Login</a></td>\n" +
		"                                            </tr>\n" +
		"                                        </table>\n" +
		"                                    </td>\n" +
		"                                </tr>\n" +
		"                            </table>\n" +
		"                        </td>\n" +
		"                    </tr> \n" +
		"    </table>\n" +
		"    <br> <br>\n" +
		"</body>" +
		"</html>"
	// message := []byte(subject + mime + body)
	return body, subject
}

func (service *AuthService) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	service.CustomLogger.InfoLogger.Info("Passwordless login confirmation with JWT token")
	token, err := jwt.ParseWithClaims(
		request.Token,
		&interceptor.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				service.CustomLogger.ErrorLogger.Error("Passwordless login confirmation with unexpected token signing method")
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return service.jwtService.publicKey, nil
		},
	)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("Passwordless login confirmation with invalid: " + request.Token)
		return nil, fmt.Errorf("Invalid token: %w", err)
	}
	user, ok := token.Claims.(*interceptor.UserClaims)
	if !ok {
		service.CustomLogger.ErrorLogger.Error("Passwordless login confirmation with invalid token claims")
		return nil, fmt.Errorf("Invalid token claims")
	}

	service.CustomLogger.SuccessLogger.Info("Passwordless login sucessfully confirmed for user with username: " + user.Username)
	return &pb.ConfirmEmailLoginResponse{
		Token: request.Token,
	}, nil
}

func (service *AuthService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	service.CustomLogger.InfoLogger.Info("User registration with username: " + request.Username + " and email: " + request.Email)
	err := checkUsernameCriteria(request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"email": request.Email,
		}).Error(err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	err = checkEmailCriteria(request.Email)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error(err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	userRequest := &user.User{
		Name:         request.Name,
		LastName:     request.LastName,
		MobileNumber: request.MobileNumber,
		Gender:       user.User_GenderEnum(request.Gender),
		Birthday:     request.Birthday,
		Username:     request.Username, // TM
		Email:        request.Email,
		Biography:    request.Biography,
		IsPublic:     false,
		IsActive:     false,
		Role:         request.Role,
	}

	for _, education := range request.Education {
		ed_id := primitive.NewObjectID().Hex()
		userRequest.Education = append(userRequest.Education, &user.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     user.Education_EducationEnum(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate,
			EndDate:   education.EndDate,
		})
	}

	for _, experience := range request.Experience {
		ex_id := primitive.NewObjectID().Hex()
		userRequest.Experience = append(userRequest.Experience, &user.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate,
			EndDate:   experience.EndDate,
		})
	}

	for _, skill := range request.Skills {
		s_id := primitive.NewObjectID().Hex()
		userRequest.Skills = append(userRequest.Skills, &user.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	for _, interest := range request.Interests {
		in_id := primitive.NewObjectID().Hex()
		userRequest.Interests = append(userRequest.Interests, &user.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	createUserRequest := &user.InsertRequest{
		User: userRequest,
	}

	err = checkPasswordCriteria(request.Password, request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"email":    userRequest.Email,
			"username": userRequest.Username,
		}).Error(err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	var authRoles []domain.Role
	for _, authRole := range request.Role {
		roles, err := service.store.FindRoleByName(authRole)
		if err != nil {
			service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
				"username": request.Username,
			}).Error("No role with name: " + authRole + " found")
			fmt.Println("Error finding role by name")
			return nil, err
		}
		authRoles = append(authRoles, *roles...)
	}

	uniqueUsername, err := service.store.IsUsernameUnique(request.Username)
	if err != nil || uniqueUsername == false {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": userRequest.Username,
		}).Error("User registration failed because username: " + request.Username + " is not unique")
		return nil, err
	}

	service.CustomLogger.DebugLogger.WithFields(logrus.Fields{
		"username": request.Username,
		"email":    request.Email,
	}).Info("Creating user")
	createUserResponse, err := service.userServiceClient.Insert(context.TODO(), createUserRequest)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": userRequest.Username,
		}).Error("Registration failed for user with username: " + request.Username)
		return nil, err
	}

	authCredentials, err := domain.NewAuthCredentials(
		createUserResponse.Id,
		request.Username,
		request.Password,
		&authRoles,
	)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("The password is not hashed for user with ID:" + createUserResponse.Id + " are not preserved")
		return nil, err
	}

	service.CustomLogger.DebugLogger.WithFields(logrus.Fields{
		"username": authCredentials.Username,
		"password": authCredentials.Password,
	}).Info("Saving user credentials with ID: " + authCredentials.Id)
	authCredentials, err = service.store.Create(authCredentials)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("Authentication credentials for user with ID:" + createUserResponse.Id + " are not created")
		return nil, err
	}

	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": userRequest.Username,
		}).Error("JWT token is not generated to user with ID:" + createUserResponse.Id)
		return nil, err
	}

	service.CustomLogger.DebugLogger.WithFields(logrus.Fields{
		"username": request.Username,
		"email":    request.Email,
	}).Info("Sending verification email to user with ID: " + authCredentials.Id)

	message, subject := verificationMailMessage(token)
	errSendingMail := sendEmail(request.Email, message, subject)
	if errSendingMail != nil {
		fmt.Println("err:  ", errSendingMail)
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": userRequest.Username,
			"email":    userRequest.Email,
		}).Error("No email was sent to user with ID:" + createUserResponse.Id)
		return nil, errSendingMail
	}

	service.CustomLogger.SuccessLogger.WithFields(logrus.Fields{
		"username": userRequest.Username,
		"email":    userRequest.Email,
	}).Info("Account successfully created for user with ID: " + createUserResponse.Id)
	return &pb.RegisterResponse{
		StatusCode: "200",
		Message:    "Success! Check your email to activate your account",
	}, nil
}

func sendEmail(sendTo, body, subject string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.NewConfig().EmailFrom)
	msg.SetHeader("To", sendTo)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	n := gomail.NewDialer(config.NewConfig().EmailHost, config.NewConfig().EmailPort, config.NewConfig().EmailFrom, config.NewConfig().EmailPassword)
	err := n.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func checkPasswordCriteria(password, username string) error {
	var err error
	var passLowercase, passUppercase, passNumber, passSpecial, passLength, passNoSpaces, passNoUsername bool
	passNoSpaces = true
	if len(password) >= minPasswordLength {
		passLength = true
	}
	if !strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		passNoUsername = true
	}
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			passLowercase = true
		case unicode.IsUpper(char):
			passUppercase = true
		case unicode.IsNumber(char):
			passNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			passSpecial = true
		case unicode.IsSpace(int32(char)):
			passNoSpaces = false
		}
	}
	if !passLowercase || !passUppercase || !passNumber || !passSpecial || !passLength || !passNoSpaces || !passNoUsername {
		switch false {
		case passLowercase:
			err = errors.New("Password must contain at least one lowercase letter")
		case passUppercase:
			err = errors.New("Password must contain at least one uppercase letter")
		case passNumber:
			err = errors.New("Password must contain at least one number")
		case passSpecial:
			err = errors.New("Password must contain at least one special character")
		case passLength:
			err = errors.New("Password must be longer than 8 characters")
		case passNoSpaces:
			err = errors.New("Password should not contain any spaces")
		case passNoUsername:
			err = errors.New("Password should not contain your username")
		}
		return err
	}
	return nil
}

func checkEmailCriteria(email string) error {
	if len(email) == 0 {
		return errors.New("Email should not be empty")
	}
	_, err := mail.ParseAddress(email)

	if err != nil {
		return errors.New("Email is invalid.")
	}
	return nil
}

func checkUsernameCriteria(username string) error {
	if len(username) == 0 {
		return errors.New("Username should not be empty")
	}

	for _, char := range username {

		if unicode.IsSpace(int32(char)) {
			return errors.New("Username should not contain any spaces")
		}
	}
	return nil
}

func (service *AuthService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	p, _ := peer.FromContext(ctx)
	service.CustomLogger.InfoLogger.WithFields(logrus.Fields{
		"ip_address": p.Addr.String(),
	}).Info("Login to application with username: " + request.Username)
	err := checkUsernameCriteria(request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("No auth credentials found with username: " + request.Username)
		fmt.Println(err.Error())
		return nil, err
	}

	err = checkPasswordCriteria(request.Password, request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error(err.Error())
		fmt.Println(err.Error())
		return nil, err
	}

	authCredentials, err := service.store.FindByUsername(request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("No auth credentials found with username: " + request.Username)
		return nil, err
	}

	service.CustomLogger.DebugLogger.Info("Getting all roles for user with ID:" + authCredentials.Id)
	var authRoles []domain.Role
	for _, authRole := range *authCredentials.Roles {
		roles, err := service.store.FindRoleByName(authRole.Name)
		if err != nil {
			service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
				"username": request.Username,
			}).Error("No role found by name: " + authRole.Name)
			fmt.Println("Error finding role by name")
			return nil, err
		}
		authRoles = append(authRoles, *roles...)
	}
	authCredentials.Roles = &authRoles

	userReq := &user.GetRequest{
		Id: authCredentials.Id,
	}
	user, err := service.userServiceClient.GetIsActive(ctx, userReq)
	if err != nil {
		fmt.Println("Error finging user data")
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("Not found user with ID: " + authCredentials.Id)
		return nil, err
	}
	if !user.IsActive {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("Not activated user with ID: " + authCredentials.Id)
		return nil, errors.New("Account is not activated")
	}

	ok := authCredentials.CheckPassword(request.Password)
	if !ok {
		service.CustomLogger.WarningLogger.WithFields(logrus.Fields{
			"username":   request.Username,
			"ip_address": p.Addr.String(),
		}).Warn("User with ID: " + authCredentials.Id + " tried to log in with the wrong credentials")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
	}

	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username":   request.Username,
			"ip_address": p.Addr.String(),
		}).Error("JWT token is not generated for user with ID: " + authCredentials.Id)
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}
	service.CustomLogger.SuccessLogger.Info("Successful user login with username: " + authCredentials.Username)
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (service *AuthService) CreateNewAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	service.CustomLogger.InfoLogger.Info("Generating API token for user: " + request.Username)
	authCredentials, err := service.store.FindByUsername(request.Username)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("No found authentication credentials with username: " + request.Username)
		return nil, err
	}

	token, hashedToken, err := service.apiTokenService.GenerateAPIToken(authCredentials)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("API token not generated for user with username: " + request.Username)
		return nil, status.Errorf(codes.Internal, "Could not generate API token")
	}

	updateCodeErr := service.store.UpdateAPIToken(authCredentials.Id, hashedToken)
	if updateCodeErr != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Error("API token not updated by user with username: " + request.Username)
		fmt.Println("Updating api token error")
		return nil, updateCodeErr
	}

	service.CustomLogger.SuccessLogger.Info("API token successfully generated for user with ID: " + authCredentials.Id)
	return &pb.NewAPITokenResponse{
		Token: token,
	}, nil
}

func (service *AuthService) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetAllResponse, error) {
	service.CustomLogger.InfoLogger.Info("Finding all auth credentials")
	auths, err := service.store.FindAll()
	if err != nil || *auths == nil {
		service.CustomLogger.ErrorLogger.Error("Error finding all auth credentials")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Auth: []*pb.Auth{},
	}

	for _, auth := range *auths {
		current := &pb.Auth{
			Id:               auth.Id,
			Username:         auth.Username,
			Password:         auth.Password,
			VerificationCode: auth.VerificationCode,
			ExpirationTime:   auth.ExpirationTime,
		}

		for _, role := range *auth.Roles {
			rolePermissions, err := service.store.GetAllPermissionsByRole(role.Name)
			if err != nil {
				service.CustomLogger.ErrorLogger.Error("Error finding all permission by role name: " + role.Name)
				fmt.Println("Greska GetAll - GetAllPermissionsByRole")
			}

			var rolePermissionsPb []*pb.Permission
			for _, perm := range *rolePermissions {
				permPb := pb.Permission{
					ID:   uint32(perm.ID),
					Name: perm.Name,
				}
				rolePermissionsPb = append(rolePermissionsPb, &permPb)
			}
			current.Roles = append(current.Roles, &pb.Role{
				ID:          uint32(role.ID),
				Name:        role.Name,
				Permissions: rolePermissionsPb,
			})
		}
		response.Auth = append(response.Auth, current)
	}
	service.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(*auths)) + " auth credentials")
	return response, nil
}

func (service *AuthService) UpdateUsername(ctx context.Context, request *pb.UpdateUsernameRequest) (*pb.UpdateUsernameResponse, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	service.CustomLogger.InfoLogger.Info("User with ID:" + userId + " is updating username")
	if userId == "" {
		return &pb.UpdateUsernameResponse{
			StatusCode: "500",
			Message:    "User id not found",
		}, nil
	} else {
		isUniqueUsername, err := service.store.IsUsernameUnique(request.Username)
		if err != nil || isUniqueUsername == false {
			service.CustomLogger.ErrorLogger.Error("User with ID:" + userId + " tried to update a non-unique username")
			return &pb.UpdateUsernameResponse{
				StatusCode: "500",
				Message:    "Username is not unique",
			}, errors.New("Username is not unique")
		}

		_, err = service.store.UpdateUsername(userId, request.Username)
		if err != nil {
			service.CustomLogger.ErrorLogger.Error("User with ID:" + userId + " did not update username")
			return &pb.UpdateUsernameResponse{
				StatusCode: "500",
				Message:    "Auth service credentials not found from JWT token",
			}, err
		}

		currentUser, err := service.userServiceClient.Get(ctx, &user.GetRequest{Id: userId})
		if err != nil {
			service.CustomLogger.ErrorLogger.Error("There is no user with with ID:" + userId)
			return nil, err
		}
		currentUser.User.Username = request.Username
		_, err = service.userServiceClient.Update(ctx, &user.UpdateRequest{User: currentUser.User})
		if err != nil {
			service.CustomLogger.ErrorLogger.Error("User with ID:" + userId + " failed to update profile")
			return nil, err
		}

		service.CustomLogger.SuccessLogger.Info("User with ID:" + userId + " has successfully updated the username")
		return &pb.UpdateUsernameResponse{
			StatusCode: "200",
			Message:    "Username updated",
		}, nil
	}
}

func (service *AuthService) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	authId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	service.CustomLogger.InfoLogger.Info("User with ID:" + authId + " is changing password")

	auth, err := service.store.FindById(authId)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("There is no auth credentials with with ID:" + authId)
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "Auth credentials not found",
		}, errors.New("Auth credentials not found")
	}

	if request.NewPassword != request.NewReenteredPassword {
		service.CustomLogger.WarningLogger.Warn("User wiht ID:" + authId + " entered passwords that do not match")
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	oldMatched := auth.CheckPassword(request.OldPassword)
	if !oldMatched {
		service.CustomLogger.WarningLogger.Warn("User wiht ID:" + authId + " has entered a password that does not match the old one")
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "Old password does not match",
		}, errors.New("Old password does not match")
	}

	err = checkPasswordCriteria(request.NewPassword, auth.Username)
	if err != nil {
		service.CustomLogger.WarningLogger.Warn("User wiht ID:" + authId + " has entered a password that does not match the old one")
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil || hashedPassword == "" {
		service.CustomLogger.WarningLogger.Warn("User wiht ID:" + authId + " has entered passwords that do not match the criteria")
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	err = service.store.UpdatePassword(authId, hashedPassword)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("User wiht ID:" + authId + " did not update the password")
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}
	service.CustomLogger.SuccessLogger.Info("User wiht ID:" + authId + " successfully updated the password")
	return &pb.ChangePasswordResponse{
		StatusCode: "200",
		Message:    "New password updated",
	}, nil
}

func verificationMailMessage(token string) (string, string) {
	// TODO SD: port se moze izvuci iz env var - 4200
	urlRedirection := "https://localhost:" + "4200" + "/activate-account/" + token
	subject := "Account activation"
	body := "<html><body style=\"background-color: #f4f4f4; margin: 0 !important; padding: 0 !important;\">\n" +
		"    <!-- HIDDEN PREHEADER TEXT -->\n" +
		"    <div style=\"display: none; font-size: 1px; color: #fefefe; line-height: 1px; font-family: 'Lato', Helvetica, Arial, sans-serif; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden;\"> We're thrilled to have you here! Get ready to dive into your new account.\n" +
		"    </div>\n" +
		"    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\">\n" +
		"        <!-- LOGO -->\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td align=\"center\" valign=\"top\" style=\"padding: 40px 10px 40px 10px;\"> </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"center\" valign=\"top\" style=\"padding: 40px 20px 20px 20px; border-radius: 4px 4px 0px 0px; color: #111111; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 48px; font-weight: 400; letter-spacing: 4px; line-height: 48px;\">\n" +
		"                            <h1 style=\"font-size: 48px; font-weight: 400; margin: 2;\">Welcome to Dislinkt!</h1> <img src=\" https://img.icons8.com/cotton/100/000000/security-checked--v3.png\" width=\"125\" height=\"120\" style=\"display: block; border: 0px;\" />\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#f4f4f4\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\" style=\"padding: 20px 30px 40px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 25px;\">\n" +
		"                            <p style=\"margin: 0;\">First, you need to activate your account. Just press the button below.</p>\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\">\n" +
		"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                <tr>\n" +
		"                                    <td bgcolor=\"#ffffff\" align=\"center\" style=\"padding: 20px 30px 60px 30px;\">\n" +
		"                                        <table border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                            <tr>\n" +
		"                                                <td align=\"center\" style=\"border-radius: 3px;\" bgcolor=\"#FFA73B\"><a href=\"" + urlRedirection + "\" target=\"_blank\" style=\"font-size: 20px; font-family: Helvetica, Arial, sans-serif; color: #ffffff; text-decoration: none; color: #ffffff; text-decoration: none; padding: 15px 25px; border-radius: 2px; border: 1px solid #FFA73B; display: inline-block;\">Activate Account</a></td>\n" +
		"                                            </tr>\n" +
		"                                        </table>\n" +
		"                                    </td>\n" +
		"                                </tr>\n" +
		"                            </table>\n" +
		"                        </td>\n" +
		"                    </tr> \n" +
		"    </table>\n" +
		"    <br> <br>\n" +
		"</body>" +
		"</html>"
	return body, subject
}

func (service *AuthService) ActivateAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	service.CustomLogger.InfoLogger.Info("Account activation with JWT token")
	p, _ := peer.FromContext(ctx)
	token, err := jwt.ParseWithClaims(
		request.Jwt,
		&interceptor.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
					"ip_address": p.Addr.String(),
				}).Error("Unexpected JWT token signing method")
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return service.jwtService.publicKey, nil
		},
	)
	if err != nil {
		service.CustomLogger.WarningLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Warn("Activation account with invalid token: " + request.Jwt)
		return nil, fmt.Errorf("Invalid token: %w", err)
	}
	claims, ok := token.Claims.(*interceptor.UserClaims)
	if !ok {
		service.CustomLogger.WarningLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Warn("Activation account with invalid token claims")
		return nil, fmt.Errorf("Invalid token claims")
	}

	id := claims.Subject
	req := &user.ActivateAccountRequest{
		Id: id,
	}
	_, err = service.userServiceClient.UpdateIsActiveById(ctx, req)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Error("Account is not activated after successfull JWT token parsing")
		return nil, err
	}

	service.CustomLogger.SuccessLogger.Info("Account successfully activated by JWT token:" + request.Jwt)
	return &pb.ActivationResponse{
		Token: request.Jwt,
	}, nil
}

func (service *AuthService) SendRecoveryCode(ctx context.Context, request *pb.SendRecoveryCodeRequest) (*pb.SendRecoveryCodeResponse, error) {
	service.CustomLogger.InfoLogger.Info("Account recovery by user email: " + request.Email)
	userServiceRequest := &user.GetIdByEmailRequest{
		Email: request.Email,
	}
	response, err := service.userServiceClient.GetIdByEmail(ctx, userServiceRequest)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("User not found by email: " + request.Email)
		fmt.Println("User not found by this email")
		fmt.Println(err)
		return nil, err
	}

	service.CustomLogger.DebugLogger.Info("Generating verification code for account recovery by email: " + request.Email)
	randomCode := rangeIn(min6DigitNumber, max6DigitNumber)
	code := strconv.Itoa(randomCode)

	expDuration := time.Duration(verificationCodeDurationInMinutes) * time.Minute
	expDate := time.Now().Add(expDuration).Unix()

	updateCodeErr := service.store.UpdateVerifactionCode(response.Id, code)
	if updateCodeErr != nil {
		service.CustomLogger.ErrorLogger.Error("Verification code for account recovery is not updated for user with email: " + request.Email)
		fmt.Println("Updating verification code error")
		return nil, updateCodeErr
	}
	updateErr := service.store.UpdateExpirationTime(response.Id, expDate)
	if updateErr != nil {
		service.CustomLogger.ErrorLogger.Error("Expiration date for account recovery is not updated for user with email: " + request.Email)
		fmt.Println("Updating expiration time error")
		return nil, updateErr
	}

	message, body := codeVerificatioMailMessage(code)
	sendingMailErr := sendEmail(request.Email, message, body)
	if sendingMailErr != nil {
		service.CustomLogger.ErrorLogger.Error("Email for account recovery is not sent to user with email: " + request.Email)
		return nil, sendingMailErr
	}

	service.CustomLogger.SuccessLogger.Info("Email for account recovery is successfully sent to user with email:" + request.Email)
	return &pb.SendRecoveryCodeResponse{
		IdAuth: response.Id,
	}, nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func codeVerificatioMailMessage(verificationCode string) (string, string) {
	subject := "Account recovery"
	// mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body style=\"background-color: #f4f4f4; margin: 0 !important; padding: 0 !important;\">\n" +
		"    <div style=\"display: none; font-size: 1px; color: #fefefe; line-height: 1px; font-family: 'Lato', Helvetica, Arial, sans-serif; max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden;\"> We're thrilled to have you here! Get ready to dive into your new account.\n" +
		"    </div>\n" +
		"    <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\">\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td align=\"center\" valign=\"top\" style=\"padding: 40px 10px 40px 10px;\"> </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#FFA73B\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"center\" valign=\"top\" style=\"padding: 40px 20px 20px 20px; border-radius: 4px 4px 0px 0px; color: #111111; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 48px; font-weight: 400; letter-spacing: 4px; line-height: 48px;\">\n" +
		"                            <h1 style=\"font-size: 48px; font-weight: 400; margin: 2;\">Verify your account</h1> <img src=\"https://img.icons8.com/external-inipagistudio-lineal-color-inipagistudio/100/000000/external-verification-email-phising-inipagistudio-lineal-color-inipagistudio.png\" width=\"125\" height=\"120\" style=\"display: block; border: 0px;\" />\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                </table>\n" +
		"            </td>\n" +
		"        </tr>\n" +
		"        <tr>\n" +
		"            <td bgcolor=\"#f4f4f4\" align=\"center\" style=\"padding: 0px 10px 0px 10px;\">\n" +
		"                <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width: 600px;\">\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\" style=\"padding: 20px 30px 40px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 25px;\">\n" +
		"                            <p style=\"margin: 0;\">To reset your password you need to verify your account with a verification code.</p>\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\">\n" +
		"                            <table width=\"100%\" border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                <tr>\n" +
		"                                    <td bgcolor=\"#ffffff\" align=\"center\" style=\"padding: 20px 30px 60px 30px;\">\n" +
		"                                        <table border=\"0\" cellspacing=\"0\" cellpadding=\"0\">\n" +
		"                                            <tr>\n" +
		"                                                <td align=\"center\" style=\"border-radius: 3px;\" >\n" +
		"                                                    <p>Your verification code:</p><h1><b> " + verificationCode + "</b></h1>\n" +
		"                                                </td>\n" +
		"                                            </tr>\n" +
		"                                        </table>\n" +
		"                                    </td>\n" +
		"                                </tr>\n" +
		"                            </table>\n" +
		"                        </td>\n" +
		"                    </tr> \n" +
		"                    <tr>\n" +
		"                        <td bgcolor=\"#ffffff\" align=\"left\" style=\"padding: 20px 30px 40px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 18px; font-weight: 400; line-height: 25px;\">\n" +
		"                            <p style=\"margin: 0;\">Sincerely,<br>Dislinkt</p>\n" +
		"                        </td>\n" +
		"                    </tr>\n" +
		"    </table>\n" +
		"    <br> <br>\n" +
		"</body>\n" +
		"</html>"
	return body, subject
}

func (service *AuthService) VerifyRecoveryCode(ctx context.Context, request *pb.VerifyRecoveryCodeRequest) (*pb.Response, error) {
	p, _ := peer.FromContext(ctx)
	service.CustomLogger.InfoLogger.WithFields(logrus.Fields{
		"ip_address": p.Addr.String(),
		"email":      request.Email,
	}).Info("Verification code for account recovery by user with ID: " + request.IdAuth)
	auth, err := service.store.FindById(request.IdAuth)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("No user found with ID: " + request.IdAuth)
		return nil, err
	}

	if auth.VerificationCode != request.VerificationCode {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
			"email":      request.Email,
		}).Error("Verification code for account recovery by user with ID: " + request.IdAuth + " is invalid")
		return &pb.Response{
			StatusCode: "500",
			Message:    "Invalid verification code",
		}, errors.New("Invalid verification code")
	}

	if auth.ExpirationTime < time.Now().Unix() {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
			"email":      request.Email,
		}).Error("Verification code for account recovery by user with ID: " + request.IdAuth + " is expired")
		return &pb.Response{
			StatusCode: "500",
			Message:    "Verification code has expired",
		}, errors.New("Verification code has expired")
	}

	updateCodeErr := service.store.UpdateVerifactionCode(request.IdAuth, "")
	if updateCodeErr != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
			"email":      request.Email,
		}).Error("Used verification code for account recovery by user with ID: " + request.IdAuth + " is not deleted")
		fmt.Println("Updating verification code error")
		return nil, updateCodeErr
	}
	updateErr := service.store.UpdateExpirationTime(request.IdAuth, 0)
	if updateErr != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
			"email":      request.Email,
		}).Error("Used verification code for account recovery by user with ID: " + request.IdAuth + " - expiration time is not updated")
		fmt.Println("Updating expiration time error")
		return nil, updateErr
	}

	service.CustomLogger.SuccessLogger.Info("Verification code for account recovery by user with ID: " + request.IdAuth + " is successfully used")
	return &pb.Response{
		StatusCode: "200",
		Message:    "Verification code is correct",
	}, nil
}

func (service *AuthService) ResetForgottenPassword(ctx context.Context, request *pb.ResetForgottenPasswordRequest) (*pb.Response, error) {
	p, _ := peer.FromContext(ctx)
	service.CustomLogger.InfoLogger.Info("User with ID: " + request.IdAuth + " recovers the forgotten password")
	auth, err := service.store.FindById(request.IdAuth)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Error("No user found with ID: " + request.IdAuth)
		return &pb.Response{
			StatusCode: "500",
			Message:    "Auth credentials not found",
		}, errors.New("Auth credentials not found")
	}

	if request.Password != request.ReenteredPassword {
		service.CustomLogger.WarningLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Warn("User with ID: " + request.IdAuth + " entered passwords that do not match")
		return &pb.Response{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	err = checkPasswordCriteria(request.Password, auth.Username)
	if err != nil {
		service.CustomLogger.WarningLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Warn("User with ID: " + request.IdAuth + " entered password that do not match with password criteria")
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil || hashedPassword == "" {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Error("Password is not successfully hashed for user with ID: " + request.IdAuth)
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	err = service.store.UpdatePassword(request.IdAuth, hashedPassword)
	if err != nil {
		service.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"ip_address": p.Addr.String(),
		}).Error("Password is not successfully updated for user with ID: " + request.IdAuth)
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}
	service.CustomLogger.SuccessLogger.Info("Password updated successfully by user with ID: " + request.IdAuth)
	return &pb.Response{
		StatusCode: "200",
		Message:    "Password updated successfully",
	}, nil
}

func (service *AuthService) GetAllPermissionsByRole(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	roleName := "Admin"
	service.CustomLogger.InfoLogger.Info("Finding role permissions by role name: " + roleName)
	_, err := service.store.GetAllPermissionsByRole(roleName)
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("No permissions found by role name: " + roleName)
		return nil, err
	}
	service.CustomLogger.SuccessLogger.Info("Permission successfully found by role name: " + roleName)
	return &pb.Response{
		StatusCode: "200",
		Message:    "OK",
	}, nil

}

func (service *AuthService) AdminsEndpoint(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	service.CustomLogger.InfoLogger.Info("Admin accesses his endpoint")
	return &pb.Response{
		StatusCode: "200",
		Message:    "OK",
	}, nil
}

func CheckString(new string, old string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(old), []byte(new))
	return err == nil
}

func (service *AuthService) GetUsernameByApiToken(ctx context.Context, request *pb.GetUsernameRequest) (*pb.GetUsernameResponse, error) {
	service.CustomLogger.InfoLogger.Info("Finding usernmae by API token")
	all, err := service.store.FindAll()
	if err != nil {
		service.CustomLogger.ErrorLogger.Error("No auth credentials found")
		return nil, err
	}

	for _, user := range *all {
		match := CheckString(request.ApiToken, user.APIToken)
		fmt.Println("BROJ 1: ", request.ApiToken)
		fmt.Println("BROJ 2: ", user.APIToken)
		fmt.Println("da li se podudaraju: ", match)

		service.CustomLogger.SuccessLogger.Info("Successfully found username: " + user.Username + " from API token")
		if match {
			return &pb.GetUsernameResponse{
				Username: user.Username,
			}, nil
		}
	}
	service.CustomLogger.ErrorLogger.Error("No username found by API token: " + request.ApiToken)
	return &pb.GetUsernameResponse{
		Username: "not found",
	}, err
}
