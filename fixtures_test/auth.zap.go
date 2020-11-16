package fixtures_test

import (
	"go.uber.org/zap/zapcore"
)

func (m *ProviderContact) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Id"
	enc.AddString(keyName, m.Id)

	keyName = "Provider"
	vv = m.Provider
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}
	return nil
}

func (m *ChangePhoneRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "OldPhone"
	enc.AddString(keyName, m.OldPhone)

	keyName = "NewPhone"
	enc.AddString(keyName, m.NewPhone)
	return nil
}

func (m *ChangePhoneResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *DetachSimProviderRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "Imsi"
	enc.AddString(keyName, m.Imsi)
	return nil
}

func (m *DetachSimProviderResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)
	return nil
}

func (m *CreateUserRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "OperationId"
	enc.AddString(keyName, m.OperationId)

	keyName = "FirstName"
	enc.AddString(keyName, m.FirstName)

	keyName = "LastName"
	enc.AddString(keyName, m.LastName)

	keyName = "Lang"
	enc.AddString(keyName, m.Lang)

	keyName = "Phone"
	enc.AddString(keyName, m.Phone)
	return nil
}

func (m *CreateUserResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "ProfileId"
	enc.AddString(keyName, m.ProfileId)
	return nil
}

func (m *DisconnectFacebookRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *DisconnectFacebookResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *JoinRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "App"
	enc.AddString(keyName, m.App)

	keyName = "DeviceInfo"
	if m.DeviceInfo != nil {
		vv = *m.DeviceInfo
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}

	keyName = "FbAccessToken"
	enc.AddString(keyName, m.FbAccessToken)

	keyName = "Email"
	enc.AddString(keyName, m.Email)

	keyName = "Phone"
	enc.AddString(keyName, m.Phone)

	keyName = "DeviceToken"
	enc.AddString(keyName, m.DeviceToken)

	keyName = "Token"
	enc.AddString(keyName, m.Token)

	keyName = "Pin"
	enc.AddString(keyName, m.Pin)

	keyName = "AdCreds"
	if m.AdCreds != nil {
		vv = *m.AdCreds
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}

	keyName = "SimCard"
	enc.AddBool(keyName, m.SimCard)
	return nil
}

func (m *JoinResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "SessionToken"
	enc.AddString(keyName, m.SessionToken)

	keyName = "DeviceToken"
	enc.AddString(keyName, m.DeviceToken)

	keyName = "Flags"
	if m.Flags != nil {
		vv = *m.Flags
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}

	keyName = "ProfileId"
	enc.AddString(keyName, m.ProfileId)

	keyName = "ResendTtl"
	enc.AddInt64(keyName, m.ResendTtl)

	keyName = "Email"
	enc.AddString(keyName, m.Email)

	keyName = "UserStatusBits"
	enc.AddUint64(keyName, m.UserStatusBits)

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)
	return nil
}

func (m *DeviceInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "MobileId"
	enc.AddString(keyName, m.MobileId)

	keyName = "UserAgent"
	enc.AddString(keyName, m.UserAgent)

	keyName = "Meta"
	if m.Meta != nil {
		vv = *m.Meta
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}
	return nil
}

func (m *DeviceMeta) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Locale"
	enc.AddString(keyName, m.Locale)

	keyName = "Timezone"
	enc.AddString(keyName, m.Timezone)
	return nil
}

func (m *UserFlags) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "EmailVerified"
	enc.AddBool(keyName, m.EmailVerified)

	keyName = "PhoneVerified"
	enc.AddBool(keyName, m.PhoneVerified)

	keyName = "ProfileCreated"
	enc.AddBool(keyName, m.ProfileCreated)

	keyName = "PinCreated"
	enc.AddBool(keyName, m.PinCreated)

	keyName = "FacebookConnected"
	enc.AddBool(keyName, m.FacebookConnected)

	keyName = "PinReset"
	enc.AddBool(keyName, m.PinReset)

	keyName = "PinRequired"
	enc.AddBool(keyName, m.PinRequired)
	return nil
}

func (m *LogoutRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "DeviceToken"
	enc.AddString(keyName, m.DeviceToken)
	return nil
}

func (m *LogoutResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "SessionToken"
	enc.AddString(keyName, m.SessionToken)
	return nil
}

func (m *EmailToken) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "Pin"
	enc.AddString(keyName, m.Pin)

	keyName = "Ttl"
	enc.AddInt64(keyName, m.Ttl)

	keyName = "Email"
	enc.AddString(keyName, m.Email)
	return nil
}

func (m *SessionToken) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "SessionId"
	enc.AddString(keyName, m.SessionId)

	keyName = "Expired"
	enc.AddInt64(keyName, m.Expired)

	keyName = "App"
	enc.AddString(keyName, m.App)

	keyName = "TypeBits"
	enc.AddUint64(keyName, m.TypeBits)

	keyName = "DeviceId"
	enc.AddString(keyName, m.DeviceId)

	keyName = "ProfileId"
	enc.AddString(keyName, m.ProfileId)

	keyName = "Locale"
	enc.AddString(keyName, m.Locale)

	keyName = "Provider"
	vv = m.Provider
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Timezone"
	enc.AddString(keyName, m.Timezone)
	return nil
}

func (m *DeviceToken) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "TokenId"
	enc.AddString(keyName, m.TokenId)

	keyName = "DeviceId"
	enc.AddString(keyName, m.DeviceId)

	keyName = "Expired"
	enc.AddInt64(keyName, m.Expired)

	keyName = "Provider"
	vv = m.Provider
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Timezone"
	enc.AddString(keyName, m.Timezone)
	return nil
}

func (m *ChangeLocaleRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Locale"
	enc.AddString(keyName, m.Locale)
	return nil
}

func (m *ChangeLocaleResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "SessionToken"
	enc.AddString(keyName, m.SessionToken)
	return nil
}

func (m *LinkContactRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Phone"
	enc.AddString(keyName, m.Phone)

	keyName = "Email"
	enc.AddString(keyName, m.Email)

	keyName = "FbAccessToken"
	enc.AddString(keyName, m.FbAccessToken)
	return nil
}

func (m *LinkContactResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "ResendTtl"
	enc.AddInt64(keyName, m.ResendTtl)
	return nil
}

func (m *VerifyContactRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Token"
	enc.AddString(keyName, m.Token)

	keyName = "Pin"
	enc.AddString(keyName, m.Pin)
	return nil
}

func (m *VerifyContactResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *UnlinkContactRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Phone"
	enc.AddString(keyName, m.Phone)

	keyName = "Email"
	enc.AddString(keyName, m.Email)

	keyName = "Facebook"
	enc.AddString(keyName, m.Facebook)
	return nil
}

func (m *UnlinkContactResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *Credentials) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Login"
	enc.AddString(keyName, m.Login)

	keyName = "Password"
	enc.AddString(keyName, m.Password)
	return nil
}

func (m *DeleteUserRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)
	return nil
}

func (m *DeleteUserResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

func (m *RegisterUserRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "SlaveId"
	enc.AddString(keyName, m.SlaveId)

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)

	keyName = "Contact"
	if m.Contact != nil {
		vv = *m.Contact
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}
	return nil
}

func (m *Contact) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "Type"
	vv = m.Type
	if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
		_ = enc.AddObject(keyName, marshaler)
	}

	keyName = "Value"
	enc.AddString(keyName, m.Value)

	keyName = "IsVerified"
	enc.AddBool(keyName, m.IsVerified)
	return nil
}

func (m *RegisterUserResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "UserId"
	enc.AddString(keyName, m.UserId)
	return nil
}

func (m *MergeUsersRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "ToUserId"
	enc.AddString(keyName, m.ToUserId)

	keyName = "FromUserId"
	enc.AddString(keyName, m.FromUserId)
	return nil
}

func (m *MergeUsersResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "MergedUserPhones"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.MergedUserPhones {
			aenc.AppendString(value)
		}
		return nil
	}))

	keyName = "MergedUserEmails"
	_ = enc.AddArray(keyName, zapcore.ArrayMarshalerFunc(func(aenc zapcore.ArrayEncoder) error {
		for _, value := range m.MergedUserEmails {
			aenc.AppendString(value)
		}
		return nil
	}))

	keyName = "MergedUserFacebook"
	enc.AddString(keyName, m.MergedUserFacebook)
	return nil
}

func (m *authClient) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var keyName string
	var vv interface{}
	_ = vv

	keyName = "cc"
	if m.cc != nil {
		vv = *m.cc
		if marshaler, ok := vv.(zapcore.ObjectMarshaler); ok {
			_ = enc.AddObject(keyName, marshaler)
		}
	}
	return nil
}

func (m *UnimplementedAuthServer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}
