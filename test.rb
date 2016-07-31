require 'sinatra'

get '/' do
  request.env.inspect
end
