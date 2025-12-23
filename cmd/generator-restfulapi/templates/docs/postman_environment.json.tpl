{
  "name": "{{.ProjectName}} - {{.Environment | title}}",
  "values": [
    {
      "key": "baseUrl",
      "value": "http://localhost:{{.ServerPort}}",
      "type": "default",
      "enabled": true
    },
    {
      "key": "apiVersion",
      "value": "v1",
      "type": "default",
      "enabled": true
    }
{{- if .EnableAuth}},
    {
      "key": "accessToken",
      "value": "",
      "type": "secret",
      "enabled": true
    },
    {
      "key": "refreshToken",
      "value": "",
      "type": "secret",
      "enabled": true
    },
    {
      "key": "testUserEmail",
      "value": "test@example.com",
      "type": "default",
      "enabled": true
    },
    {
      "key": "testUserPassword",
      "value": "password123",
      "type": "secret",
      "enabled": true
    }
{{- end}},
    {
      "key": "dbHost",
      "value": "{{.DBHost}}",
      "type": "default",
      "enabled": true
    },
    {
      "key": "dbPort",
      "value": "{{.DBPort}}",
      "type": "default",
      "enabled": true
    },
    {
      "key": "dbName",
      "value": "{{.DBName}}",
      "type": "default",
      "enabled": true
    }
  ],
  "_postman_variable_scope": "environment"
}
