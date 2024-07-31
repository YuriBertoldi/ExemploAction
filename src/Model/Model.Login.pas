unit Model.Login;

interface

uses Model.Usuario;

type
  TLoginResult = (lrSuccess, lrInvalidCredentials, lrServerError);

  TLoginResultHelper = Record helper for TLoginResult
    function Text : string;
  end;

  TLogin = class
  public
    function AuthenticateUser(const User: TUser): TLoginResult;
  end;

implementation

function TLogin.AuthenticateUser(const User: TUser): TLoginResult;
begin
  if (User.Username = 'user' ) and (User.Password = 'password') then
    Result := lrSuccess
  else
    Result := lrInvalidCredentials;
end;

{ TLoginResultHelper }

function TLoginResultHelper.Text: string;
begin
  case self of
    lrSuccess           : Result := 'Authenticated successfully.';
    lrInvalidCredentials: Result := 'Incorrect username or password.';
    lrServerError       : Result := 'Server error.';
  end;
end;


end.

