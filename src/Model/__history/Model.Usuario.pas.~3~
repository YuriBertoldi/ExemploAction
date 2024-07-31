unit Model.Usuario;

interface

uses
  System.SysUtils;

type
  TUsuario = class
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

constructor TUsuario.Create;
begin
  FUsername := EmptyStr;
  FPassword := EmptyStr;
end;

procedure TUsuario.SetPassword(const Value: string);
begin
  FPassword := Value;
end;

procedure TUsuario.SetUsername(const Value: string);
begin
  FUsername := Value;
end;

end.

