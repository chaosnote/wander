#!bin/bash

project_dir=/app/wander/env_local

## build dc
cd $project_dir/data_center
go build -o /dist/dc . 

## build game_0000
cd $project_dir/game_0000
go build -o /dist/game_0000 . 
