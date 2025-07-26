export const statusColorMapping = {
  not_checked: 'warning.light',
  checked: 'success.light',
  partially_checked: 'info.light',
  failed_check: 'error.light',
  default: 'grey.700',
};

export const statusNameMapping = {
  not_checked: 'Не проверено',
  checked: 'Проверено',
  partially_checked: 'Проверено частично',
  failed_check: 'Провалено',
};

export const getStatusColor = (status) => statusColorMapping[status] || statusColorMapping.default;
export const getStatusName = (status) => statusNameMapping[status] || status;