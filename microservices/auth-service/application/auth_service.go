package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/startup/config"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var verificationCodeDurationInMinutes int = 5
var min6DigitNumber int = 100000
var max6DigitNumber int = 999999

type AuthService struct {
	store             *persistence.AuthPostgresStore
	jwtService        *JWTService
	userServiceClient user.UserServiceClient
}

func NewAuthService(store *persistence.AuthPostgresStore, jwtService *JWTService, userServiceClient user.UserServiceClient) *AuthService {
	return &AuthService{
		store:             store,
		jwtService:        jwtService,
		userServiceClient: userServiceClient,
	}
}

func (service *AuthService) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {

	getUserRequest := &user.GetIdByEmailRequest{
		Email: request.Email,
	}

	user, err := service.userServiceClient.GetIdByEmail(context.TODO(), getUserRequest)

	if err != nil {
		return nil, errors.New("there is no user with that email or account is not activated")
	}

	authCredentials, err := service.store.FindById(user.Id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	from := config.NewConfig().EmailFrom
	password := config.NewConfig().EmailPassword

	to := []string{
		request.Email,
	}

	smtpHost := config.NewConfig().EmailHost
	smtpPort := config.NewConfig().EmailPort

	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}

	message := passwordlessLoginMailMessage(token)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error while sending mail")
	}

	return &pb.PasswordlessLoginResponse{
		Success: "Email Sent Successfully! Check your email.",
	}, nil
}

func passwordlessLoginMailMessage(token string) []byte {
	urlRedirection := "http://" + config.NewConfig().FrontendHost + ":" + config.NewConfig().FrontendPort + "/confirmed-mail/" + token

	subject := "Subject: Passwordless login\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
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
	message := []byte(subject + mime + body)
	return message
}

func (service *AuthService) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {

	token, err := jwt.ParseWithClaims(
		request.Token,
		&interceptor.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return service.jwtService.publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}
	_, ok := token.Claims.(*interceptor.UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	return &pb.ConfirmEmailLoginResponse{
		Token: request.Token,
	}, nil
}

func (service *AuthService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userRequest := &user.User{
		Name:         request.Name,
		LastName:     request.LastName,
		MobileNumber: request.MobileNumber,
		Gender:       user.User_GenderEnum(request.Gender),
		Birthday:     request.Birthday,
		Email:        request.Email,
		Biography:    request.Biography,
		IsPublic:     request.IsPublic,
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

	err := checkPasswordCriteria(request.Password, request.Username)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	createUserResponse, err := service.userServiceClient.Insert(context.TODO(), createUserRequest)
	if err != nil {
		return nil, err
	}

	authCredentials, err := domain.NewAuthCredentials(
		createUserResponse.Id,
		request.Username,
		request.Password,
		request.Role,
	)
	if err != nil {
		return nil, err
	}

	authCredentials, err = service.store.Create(authCredentials)
	if err != nil {
		return nil, err
	}

	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		return nil, err
	}
	message := verificationMailMessage(token)
	errSendingMail := sendMail(request.Email, message)
	if errSendingMail != nil {
		fmt.Println("err:  ", errSendingMail)
		return nil, errSendingMail
	}

	return &pb.RegisterResponse{
		StatusCode: "200",
		Message:    "Success! Check your email to activate your account",
	}, nil
}

func checkPasswordCriteria(password, username string) error {
	var err error
	var passLowercase, passUppercase, passNumber, passSpecial, passLength, passNoSpaces, passNoUsername bool
	passNoSpaces = true
	if len(password) >= 8 {
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

func (service *AuthService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	authCredentials, err := service.store.FindByUsername(request.Username)
	if err != nil {
		return nil, err
	}
	fmt.Println("No error finding auth credentials")

	userReq := &user.GetRequest{
		Id: authCredentials.Id,
	}
	user, err := service.userServiceClient.GetIsActive(ctx, userReq)
	if err != nil {
		fmt.Println("Error finging user data")
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("Account is not activated")
	}

	ok := authCredentials.CheckPassword(request.Password)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
	}
	fmt.Println("No error validating password")
	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (service *AuthService) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetAllResponse, error) {
	auths, err := service.store.FindAll()
	if err != nil || *auths == nil {
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
			Role:             auth.Role,
			VerificationCode: auth.VerificationCode,
			ExpirationTime:   auth.ExpirationTime,
		}
		response.Auth = append(response.Auth, current)
	}
	return response, nil
}

func (service *AuthService) UpdateUsername(ctx context.Context, request *pb.UpdateUsernameRequest) (*pb.UpdateUsernameResponse, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	if userId == "" {
		return &pb.UpdateUsernameResponse{
			StatusCode: "500",
			Message:    "User id not found",
		}, nil
	} else {
		auths, err := service.store.FindAll()
		for _, auth := range *auths {
			if auth.Username == request.Username {
				log.Println("Username is not unique")
				return &pb.UpdateUsernameResponse{
					StatusCode: "500",
					Message:    "Username is not unique",
				}, errors.New("Username is not unique")
			}
		}
		response, err := service.store.UpdateUsername(userId, request.Username)
		if err != nil {
			return &pb.UpdateUsernameResponse{
				StatusCode: "500",
				Message:    "Auth service credentials not found from JWT token",
			}, err
		}
		log.Print(response)
		return &pb.UpdateUsernameResponse{
			StatusCode: "200",
			Message:    "Username updated",
		}, nil
	}
}

func (service *AuthService) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	authId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	auth, err := service.store.FindById(authId)
	if err != nil {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "Auth credentials not found",
		}, errors.New("Auth credentials not found")
	}

	if request.NewPassword != request.NewReenteredPassword {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	oldMatched := auth.CheckPassword(request.OldPassword)
	if !oldMatched {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    "Old password does not match",
		}, errors.New("Old password does not match")
	}

	err = checkPasswordCriteria(request.NewPassword, auth.Username)
	if err != nil {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	hashedPassword, err := auth.HashPassword(request.NewPassword)
	if err != nil || hashedPassword == "" {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	err = service.store.UpdatePassword(authId, hashedPassword)
	if err != nil {
		return &pb.ChangePasswordResponse{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}
	return &pb.ChangePasswordResponse{
		StatusCode: "200",
		Message:    "New password updated",
	}, nil
}

func sendMail(emailTo string, message []byte) error {
	from := config.NewConfig().EmailFrom
	emailPassword := config.NewConfig().EmailPassword
	to := []string{emailTo}

	host := config.NewConfig().EmailHost
	port := config.NewConfig().EmailPort
	smtpAddress := host + ":" + port

	authMail := smtp.PlainAuth("", from, emailPassword, host)

	errSendingMail := smtp.SendMail(smtpAddress, authMail, from, to, message)
	if errSendingMail != nil {
		fmt.Println("err:  ", errSendingMail)
		return errSendingMail
	}
	return nil
}

func verificationMailMessage(token string) []byte {
	// TODO SD: port se moze izvuci iz env var - 4200
	urlRedirection := "http://localhost:" + "4200" + "/activate-account/" + token

	subject := "Subject: Account activation\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
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
	message := []byte(subject + mime + body)
	return message
}

func (service *AuthService) ActivateAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	token, err := jwt.ParseWithClaims(
		request.Jwt,
		&interceptor.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return service.jwtService.publicKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}
	claims, ok := token.Claims.(*interceptor.UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	id := claims.Subject
	req := &user.ActivateAccountRequest{
		Id: id,
	}
	_, err = service.userServiceClient.UpdateIsActiveById(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.ActivationResponse{
		Token: request.Jwt,
	}, nil
}

func (service *AuthService) SendRecoveryCode(ctx context.Context, request *pb.SendRecoveryCodeRequest) (*pb.SendRecoveryCodeResponse, error) {
	userServiceRequest := &user.GetIdByEmailRequest{
		Email: request.Email,
	}
	response, err := service.userServiceClient.GetIdByEmail(ctx, userServiceRequest)
	if err != nil {
		fmt.Println("User not found by this email")
		fmt.Println(err)
		return nil, err
	}

	randomCode := rangeIn(min6DigitNumber, max6DigitNumber)
	fmt.Println(randomCode)
	code := strconv.Itoa(randomCode)

	expDuration := time.Duration(verificationCodeDurationInMinutes) * time.Minute
	expDate := time.Now().Add(expDuration).Unix()

	updateCodeErr := service.store.UpdateVerifactionCode(response.Id, code)
	if updateCodeErr != nil {
		fmt.Println("Updating verification code error")
		return nil, updateCodeErr
	}
	updateErr := service.store.UpdateExpirationTime(response.Id, expDate)
	if updateErr != nil {
		fmt.Println("Updating expiration time error")
		return nil, updateErr
	}

	message := codeVerificatioMailMessage(code)
	sendingMailErr := sendMail(request.Email, message)
	if sendingMailErr != nil {
		return nil, sendingMailErr
	}

	return &pb.SendRecoveryCodeResponse{
		IdAuth: response.Id,
	}, nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func codeVerificatioMailMessage(verificationCode string) []byte {
	subject := "Subject: Account recovery\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
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
	message := []byte(subject + mime + body)
	return message
}

func (service *AuthService) VerifyRecoveryCode(ctx context.Context, request *pb.VerifyRecoveryCodeRequest) (*pb.Response, error) {
	auth, err := service.store.FindById(request.IdAuth)
	if err != nil {
		return nil, err
	}

	if auth.VerificationCode != request.VerificationCode {
		return &pb.Response{
			StatusCode: "500",
			Message:    "Invalid verification code",
		}, errors.New("Invalid verification code")
	}

	if auth.ExpirationTime < time.Now().Unix() {
		return &pb.Response{
			StatusCode: "500",
			Message:    "Verification code has expired",
		}, errors.New("Verification code has expired")
	}

	updateCodeErr := service.store.UpdateVerifactionCode(request.IdAuth, "")
	if updateCodeErr != nil {
		fmt.Println("Updating verification code error")
		return nil, updateCodeErr
	}
	updateErr := service.store.UpdateExpirationTime(request.IdAuth, 0)
	if updateErr != nil {
		fmt.Println("Updating expiration time error")
		return nil, updateErr
	}

	return &pb.Response{
		StatusCode: "200",
		Message:    "Verification code is correct",
	}, nil
}

func (service *AuthService) ResetForgottenPassword(ctx context.Context, request *pb.ResetForgottenPasswordRequest) (*pb.Response, error) {
	auth, err := service.store.FindById(request.IdAuth)
	if err != nil {
		return &pb.Response{
			StatusCode: "500",
			Message:    "Auth credentials not found",
		}, errors.New("Auth credentials not found")
	}

	if request.Password != request.ReenteredPassword {
		return &pb.Response{
			StatusCode: "500",
			Message:    "New passwords do not match",
		}, errors.New("New passwords do not match")
	}

	err = checkPasswordCriteria(request.Password, auth.Username)
	if err != nil {
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil || hashedPassword == "" {
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	err = service.store.UpdatePassword(request.IdAuth, hashedPassword)
	if err != nil {
		return &pb.Response{
			StatusCode: "500",
			Message:    err.Error(),
		}, err
	}

	return &pb.Response{
		StatusCode: "200",
		Message:    "Password updated successfully",
	}, nil
}
