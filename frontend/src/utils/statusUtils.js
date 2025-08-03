export const statusColorMapping = {
  not_checked: 'warning',
  checked: 'success',
  partially_checked: 'info',
  failed_check: 'error',
  default: 'grey',
};

export const statusNameMapping = {
  not_checked: 'Не проверено',
  checked: 'Проверено',
  partially_checked: 'Проверено частично',
  failed_check: 'Провалено',
};

export const getStatusChipColor = (status) => {
  return statusColorMapping[status] || 'default';
};
export const getStatusName = (status) => statusNameMapping[status] || status;