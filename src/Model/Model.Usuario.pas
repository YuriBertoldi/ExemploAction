unit Model.Usuario;

interface

uses
  System.SysUtils;

type
  TUser = class
  private
    FPassword: string;
    FUsername: string;
    procedure SetPassword(const Value: string);
    procedure SetUsername(const Value: string);
  public
    constructor Create;
    property Username: string read FUsername write SetUsername;
    property Password: string read FPassword write SetPassword;
  end;

implementation

constructor TUser.Create;
begin
  FUsername := EmptyStr;
  FPassword := EmptyStr;
end;

procedure TUser.SetPassword(const Value: string);
begin
  FPassword := Value;
end;

procedure TUser.SetUsername(const Value: string);
begin
  FUsername := Value;
end;

end.

