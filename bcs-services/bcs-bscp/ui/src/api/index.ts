import http from "../request"

/**
 * 获取应用列表
 * @param biz_id 业务ID
 * @param params 查询过滤条件
 * @returns 
 */
export type IAppListQuery = {
  operator?: string,
  name?: string,
  start?: number,
  limit?: number 
}

export const getAppList = (biz_id: number, params: IAppListQuery = {}) => {
  return http.get(`config/list/app/app/biz_id/${biz_id}`, { params }).then(resp => resp.data);
}

/**
 * 删除应用
 * @param id 应用ID
 * @param biz_id 业务ID
 * @returns 
 */
export const deleteApp = (id: number, biz_id: number) => {
  return http.delete(`config/delete/app/app/app_id/${id}/biz_id/${biz_id}`);
}

/**
 * 创建应用
 * @param biz_id 业务ID
 * @param params 
 * @returns 
 */
export const createApp = (biz_id: number, params: any) => {
  return http.post(`config/create/app/app/biz_id/${biz_id}`, { biz_id, ...params }).then(resp => resp.data);
}

/**
 * 更新应用
 * @param params { id, biz_id, name?, memo?, reload_type?, reload_file_path? }
 * @returns 
 */
export const updateApp = (params: any) => {
  const { id, biz_id, data } = params;
  return http.put(`config/update/app/app/app_id/${id}/biz_id/${biz_id}`, data).then(resp => resp.data);
}