﻿unit Controller.Login;

interface

uses
  Model.login, Model.Usuario, FMX.StdCtrls;

type
  TLoginController = class
  private
    FLoginManager: TLogin;
    FUser : tUser;
    procedure UpdateMessage(TypeMessage: TLoginResult; LabelMessage : tLabel);
  public
    constructor Create;
    destructor Destroy; override;
    function DoLogin(LabelMessage : tLabel): TLoginResult;
    procedure SetUser(const Username, Password: string);
  end;

implementation

constructor TLoginController.Create;
begin
  FLoginManager := TLogin.Create;
  FUser         := tUser.Create;
end;

destructor TLoginController.Destroy;
begin
  FLoginManager.Free;
  FUser.Free;
  inherited;
end;

function TLoginController.DoLogin(LabelMessage : tLabel): TLoginResult;
begin
  Result := FLoginManager.AuthenticateUser(FUser);
  UpdateMessage(Result, LabelMessage);
end;

procedure TLoginController.SetUser(const Username, Password: string);
begin
  FUser.Username := username;
  FUser.Password := Password;
end;

procedure TLoginController.UpdateMessage(TypeMessage: TLoginResult;
  LabelMessage: tLabel);
begin
  LabelMessage.Text := TypeMessage.Text;
end;

end.

