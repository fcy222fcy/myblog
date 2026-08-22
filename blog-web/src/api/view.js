import request from './request'

export const recordContentView = (contentType, contentID) => {
  return request.post('/views', {
    content_type: contentType,
    content_id: contentID
  })
}
