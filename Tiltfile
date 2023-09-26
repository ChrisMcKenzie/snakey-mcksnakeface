local_resource(
    'battlesnake', 
    cmd='go build',
    serve_cmd='./starter-snake-go',
    serve_env={
        'PORT': "80",
    },
    deps=['main.go', './', 'models.go', 'server.go'],
    ignore=['starter-snake-go'],
)
