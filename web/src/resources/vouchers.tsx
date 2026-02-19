// web/src/resources/vouchers.tsx
import {
  List,
  Datagrid,
  TextField,
  DateField,
  Edit,
  TextInput,
  Create,
  Show,
  BooleanInput,
  NumberInput,
  required,
  useRecordContext,
  Toolbar,
  SaveButton,
  DeleteButton,
  SimpleForm,
  TopToolbar,
  ListButton,
  CreateButton,
  useTranslate,
  useListContext,
  FunctionField,
  SelectInput,
  ReferenceInput,
  useNotify,
} from 'react-admin';
import {
  Box,
  Chip,
  Typography,
  Card,
  Avatar,
  IconButton,
  Tooltip,
} from '@mui/material';
import { useCallback } from 'react';
import {
  ConfirmationNumber as VoucherIcon,
  ContentCopy as CopyIcon,
  CheckCircle as EnabledIcon,
  Cancel as DisabledIcon,
} from '@mui/icons-material';
import {
  ServerPagination,
  FormSection,
  FieldGrid,
  FieldGridItem,
  formLayoutSx,
  controlWrapperSx,
  DetailSectionCard,
  DetailItem,
} from '../components';

const LARGE_LIST_PER_PAGE = 50;

// ============ 类型定义 ============

interface VoucherBatch {
  id: number;
  name?: string;
  node_id?: number;
  profile_id?: number;
  total_count?: number;
  used_count?: number;
  expire_time?: string;
  valid_days?: number;
  prefix?: string;
  code_length?: number;
  status?: 'enabled' | 'disabled';
  remark?: string;
  created_at?: string;
}

interface Voucher {
  id: number;
  batch_id?: number;
  code?: string;
  password?: string;
  profile_id?: number;
  status?: 'available' | 'used' | 'expired' | 'disabled';
  user_id?: number;
  redeemed_at?: string;
  expire_time?: string;
  remark?: string;
  created_at?: string;
}

// ============ 工具函数 ============

const formatTimestamp = (value?: string | number): string => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString();
};

const formatDate = (value?: string | number): string => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleDateString();
};

// ============ 状态组件 ============

const BatchStatusIndicator = ({ status }: { status?: string }) => {
  const translate = useTranslate();
  const isEnabled = status === 'enabled';
  return (
    <Chip
      icon={isEnabled ? <EnabledIcon sx={{ fontSize: '0.85rem !important' }} /> : <DisabledIcon sx={{ fontSize: '0.85rem !important' }} />}
      label={isEnabled ? translate('resources.voucher-batches.status.enabled', { _: '启用' }) : translate('resources.voucher-batches.status.disabled', { _: '禁用' })}
      size="small"
      color={isEnabled ? 'success' : 'default'}
      variant={isEnabled ? 'filled' : 'outlined'}
      sx={{ height: 22, fontWeight: 500, fontSize: '0.75rem' }}
    />
  );
};

const VoucherStatusIndicator = ({ status }: { status?: string }) => {
  const translate = useTranslate();
  const statusConfig: Record<string, { color: 'success' | 'warning' | 'error' | 'default'; label: string }> = {
    available: { color: 'success', label: translate('resources.vouchers.status.available', { _: '可用' }) },
    used: { color: 'warning', label: translate('resources.vouchers.status.used', { _: '已使用' }) },
    expired: { color: 'error', label: translate('resources.vouchers.status.expired', { _: '已过期' }) },
    disabled: { color: 'default', label: translate('resources.vouchers.status.disabled', { _: '已禁用' }) },
  };
  const config = statusConfig[status || 'available'];
  return (
    <Chip
      label={config.label}
      size="small"
      color={config.color}
      variant="outlined"
      sx={{ height: 22, fontWeight: 500, fontSize: '0.75rem' }}
    />
  );
};

// ============ 字段组件 ============

const BatchNameField = () => {
  const record = useRecordContext<VoucherBatch>();
  if (!record) return null;

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <Avatar
        sx={{
          width: 32,
          height: 32,
          fontSize: '0.85rem',
          fontWeight: 600,
          bgcolor: record.status === 'enabled' ? 'primary.main' : 'grey.400',
        }}
      >
        <VoucherIcon sx={{ fontSize: 18 }} />
      </Avatar>
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 600, color: 'text.primary', lineHeight: 1.3 }}>
          {record.name || '-'}
        </Typography>
        <BatchStatusIndicator status={record.status} />
      </Box>
    </Box>
  );
};

const UsageProgressField = () => {
  const record = useRecordContext<VoucherBatch>();
  if (!record) return null;

  const used = record.used_count || 0;
  const total = record.total_count || 1;
  const percentage = Math.round((used / total) * 100);

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <Box sx={{ flex: 1 }}>
        <Box
          sx={{
            height: 6,
            borderRadius: 3,
            bgcolor: 'grey.200',
            overflow: 'hidden',
          }}
        >
          <Box
            sx={{
              height: '100%',
              width: `${percentage}%`,
              borderRadius: 3,
              bgcolor: percentage > 90 ? 'error.main' : percentage > 70 ? 'warning.main' : 'success.main',
              transition: 'width 0.3s ease',
            }}
          />
        </Box>
      </Box>
      <Typography variant="caption" sx={{ minWidth: 50, textAlign: 'right', color: 'text.secondary' }}>
        {used}/{total}
      </Typography>
    </Box>
  );
};

const VoucherCodeField = () => {
  const record = useRecordContext<Voucher>();
  const notify = useNotify();

  const handleCopy = useCallback(() => {
    if (record?.code) {
      navigator.clipboard.writeText(record.code);
      notify('Voucher code copied to clipboard', { type: 'info' });
    }
  }, [record?.code, notify]);

  if (!record) return null;

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
      <Typography
        variant="body2"
        sx={{
          fontFamily: 'monospace',
          fontWeight: 600,
          letterSpacing: 0.5,
        }}
      >
        {record.code}
      </Typography>
      <Tooltip title="Copy code">
        <IconButton size="small" onClick={handleCopy} sx={{ p: 0.25 }}>
          <CopyIcon sx={{ fontSize: 14 }} />
        </IconButton>
      </Tooltip>
    </Box>
  );
};

// ============ 列表操作栏 ============

const BatchListActions = () => {
  const translate = useTranslate();
  return (
    <TopToolbar>
      <CreateButton label={translate('resources.voucher-batches.actions.create', { _: '新建批次' })} />
    </TopToolbar>
  );
};

const VoucherListActions = () => {
  const translate = useTranslate();
  return (
    <TopToolbar>
      <ListButton
        resource="voucher-batches"
        label={translate('resources.vouchers.actions.batches', { _: '批次管理' })}
        icon={<VoucherIcon />}
      />
    </TopToolbar>
  );
};

// ============ 列表内容 ============

const BatchListContent = () => {
  const translate = useTranslate();
  const { data, isLoading, total } = useListContext<VoucherBatch>();

  if (isLoading) {
    return <Typography>Loading...</Typography>;
  }

  if (!data || data.length === 0) {
    return (
      <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}` }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 8, color: 'text.secondary' }}>
          <VoucherIcon sx={{ fontSize: 64, opacity: 0.3, mb: 2 }} />
          <Typography variant="h6" sx={{ opacity: 0.6, mb: 1 }}>
            {translate('resources.voucher-batches.empty.title', { _: '暂无批次' })}
          </Typography>
        </Box>
      </Card>
    );
  }

  return (
    <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}`, overflow: 'hidden' }}>
      <Box sx={{ px: 2, py: 1, bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.02)' : 'rgba(0,0,0,0.01)', borderBottom: theme => `1px solid ${theme.palette.divider}` }}>
        <Typography variant="body2" color="text.secondary">
          共 <strong>{total?.toLocaleString() || 0}</strong> 个批次
        </Typography>
      </Box>
      <Box sx={{ overflowX: 'auto' }}>
        <Datagrid rowClick="show" bulkActionButtons={false}>
          <FunctionField
            source="name"
            label={translate('resources.voucher-batches.fields.name', { _: '批次名称' })}
            render={() => <BatchNameField />}
          />
          <TextField
            source="total_count"
            label={translate('resources.voucher-batches.fields.total_count', { _: '总数' })}
          />
          <FunctionField
            source="used_count"
            label={translate('resources.voucher-batches.fields.usage', { _: '使用情况' })}
            render={() => <UsageProgressField />}
          />
          <FunctionField
            source="expire_time"
            label={translate('resources.voucher-batches.fields.expire_time', { _: '过期时间' })}
            render={(record: VoucherBatch) => formatDate(record.expire_time)}
          />
          <DateField
            source="created_at"
            label={translate('resources.voucher-batches.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Box>
    </Card>
  );
};

// ============ 列表页面 ============

export const VoucherBatchList = () => {
  return (
    <List
      actions={<BatchListActions />}
      sort={{ field: 'created_at', order: 'DESC' }}
      perPage={LARGE_LIST_PER_PAGE}
      pagination={<ServerPagination />}
      empty={false}
    >
      <BatchListContent />
    </List>
  );
};

export const VoucherList = () => {
  const translate = useTranslate();
  return (
    <List
      actions={<VoucherListActions />}
      sort={{ field: 'created_at', order: 'DESC' }}
      perPage={LARGE_LIST_PER_PAGE}
      pagination={<ServerPagination />}
      empty={false}
    >
      <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}`, overflow: 'hidden' }}>
        <Datagrid bulkActionButtons={false}>
          <FunctionField
            source="code"
            label={translate('resources.vouchers.fields.code', { _: '券码' })}
            render={() => <VoucherCodeField />}
          />
          <FunctionField
            source="status"
            label={translate('resources.vouchers.fields.status', { _: '状态' })}
            render={(record: Voucher) => <VoucherStatusIndicator status={record.status} />}
          />
          <TextField
            source="batch_id"
            label={translate('resources.vouchers.fields.batch_id', { _: '批次ID' })}
          />
          <FunctionField
            source="expire_time"
            label={translate('resources.vouchers.fields.expire_time', { _: '过期时间' })}
            render={(record: Voucher) => formatDate(record.expire_time)}
          />
          <DateField
            source="created_at"
            label={translate('resources.vouchers.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Card>
    </List>
  );
};

// BatchVoucherList displays vouchers filtered by batch_id, used inside VoucherBatchShow
const BatchVoucherList = ({ batchId }: { batchId: number }) => {
  const translate = useTranslate();
  return (
    <List
      resource="vouchers"
      filter={{ batch_id: batchId }}
      sort={{ field: 'created_at', order: 'DESC' }}
      perPage={LARGE_LIST_PER_PAGE}
      pagination={<ServerPagination />}
      empty={false}
      actions={false}
    >
      <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}`, overflow: 'hidden' }}>
        <Datagrid bulkActionButtons={false}>
          <FunctionField
            source="code"
            label={translate('resources.vouchers.fields.code', { _: '券码' })}
            render={() => <VoucherCodeField />}
          />
          <FunctionField
            source="status"
            label={translate('resources.vouchers.fields.status', { _: '状态' })}
            render={(record: Voucher) => <VoucherStatusIndicator status={record.status} />}
          />
          <FunctionField
            source="expire_time"
            label={translate('resources.vouchers.fields.expire_time', { _: '过期时间' })}
            render={(record: Voucher) => formatDate(record.expire_time)}
          />
          <DateField
            source="created_at"
            label={translate('resources.vouchers.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Card>
    </List>
  );
};

// ============ 创建/编辑表单 ============

const BatchFormToolbar = (props: any) => (
  <Toolbar {...props}>
    <SaveButton />
    <DeleteButton mutationMode="pessimistic" />
  </Toolbar>
);

export const VoucherBatchCreate = () => {
  const translate = useTranslate();

  return (
    <Create>
      <SimpleForm sx={formLayoutSx}>
        <FormSection
          title={translate('resources.voucher-batches.sections.basic.title', { _: '基本信息' })}
          description={translate('resources.voucher-batches.sections.basic.description', { _: '批次的基本配置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="name"
                label={translate('resources.voucher-batches.fields.name', { _: '批次名称' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <ReferenceInput
                source="profile_id"
                reference="radius-profiles"
                label={translate('resources.voucher-batches.fields.profile_id', { _: '计费策略' })}
              >
                <SelectInput optionText="name" validate={[required()]} fullWidth size="small" />
              </ReferenceInput>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.voucher-batches.sections.voucher.title', { _: '券码配置' })}
          description={translate('resources.voucher-batches.sections.voucher.description', { _: '券码生成设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 3 }}>
            <FieldGridItem>
              <NumberInput
                source="total_count"
                label={translate('resources.voucher-batches.fields.total_count', { _: '券码数量' })}
                defaultValue={100}
                min={1}
                max={10000}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="prefix"
                label={translate('resources.voucher-batches.fields.prefix', { _: '券码前缀' })}
                helperText={translate('resources.voucher-batches.helpers.prefix', { _: '可选的前缀，最多10个字符' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="code_length"
                label={translate('resources.voucher-batches.fields.code_length', { _: '券码长度' })}
                defaultValue={10}
                min={6}
                max={32}
                helperText={translate('resources.voucher-batches.helpers.code_length', { _: '随机部分长度' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.voucher-batches.sections.validity.title', { _: '有效期配置' })}
          description={translate('resources.voucher-batches.sections.validity.description', { _: '券码有效期设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="expire_time"
                label={translate('resources.voucher-batches.fields.expire_time', { _: '批次过期时间' })}
                type="datetime-local"
                validate={[required()]}
                fullWidth
                size="small"
                InputLabelProps={{ shrink: true }}
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="valid_days"
                label={translate('resources.voucher-batches.fields.valid_days', { _: '有效天数' })}
                min={0}
                max={3650}
                helperText={translate('resources.voucher-batches.helpers.valid_days', { _: '兑换后有效天数，0表示使用批次过期时间' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.voucher-batches.sections.remark.title', { _: '备注' })}
        >
          <FieldGrid columns={{ xs: 1 }}>
            <FieldGridItem>
              <TextInput
                source="remark"
                label={translate('resources.voucher-batches.fields.remark', { _: '备注' })}
                multiline
                minRows={3}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Create>
  );
};

export const VoucherBatchEdit = () => {
  const translate = useTranslate();

  return (
    <Edit>
      <SimpleForm toolbar={<BatchFormToolbar />} sx={formLayoutSx}>
        <FormSection
          title={translate('resources.voucher-batches.sections.basic.title', { _: '基本信息' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput source="id" disabled fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="name"
                label={translate('resources.voucher-batches.fields.name', { _: '批次名称' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput
                  source="status"
                  label={translate('resources.voucher-batches.fields.status_enabled', { _: '启用状态' })}
                />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.voucher-batches.sections.validity.title', { _: '有效期配置' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="expire_time"
                label={translate('resources.voucher-batches.fields.expire_time', { _: '过期时间' })}
                type="datetime-local"
                fullWidth
                size="small"
                InputLabelProps={{ shrink: true }}
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="valid_days"
                label={translate('resources.voucher-batches.fields.valid_days', { _: '有效天数' })}
                min={0}
                max={3650}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.voucher-batches.sections.remark.title', { _: '备注' })}>
          <FieldGrid columns={{ xs: 1 }}>
            <FieldGridItem>
              <TextInput source="remark" multiline minRows={3} fullWidth size="small" />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Edit>
  );
};

// ============ 详情页 ============

// VoucherBatchDetails is the inner component that uses useRecordContext inside Show
const VoucherBatchDetails = () => {
  const record = useRecordContext<VoucherBatch>();
  const translate = useTranslate();

  if (!record) return null;

  return (
    <Box sx={{ p: 3 }}>
      <DetailSectionCard
        title={translate('resources.voucher-batches.sections.info.title', { _: '批次信息' })}
        icon={<VoucherIcon />}
      >
        <Box sx={{ display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)' } }}>
          <DetailItem label={translate('resources.voucher-batches.fields.id', { _: 'ID' })} value={String(record.id || '-')} />
          <DetailItem label={translate('resources.voucher-batches.fields.name', { _: '名称' })} value={record.name || '-'} />
          <DetailItem label={translate('resources.voucher-batches.fields.total_count', { _: '总数' })} value={String(record.total_count || 0)} />
          <DetailItem label={translate('resources.voucher-batches.fields.used_count', { _: '已使用' })} value={String(record.used_count || 0)} />
          <DetailItem label={translate('resources.voucher-batches.fields.expire_time', { _: '过期时间' })} value={formatTimestamp(record.expire_time)} />
          <DetailItem label={translate('resources.voucher-batches.fields.status', { _: '状态' })} value={<BatchStatusIndicator status={record.status} />} />
          <DetailItem label={translate('resources.voucher-batches.fields.created_at', { _: '创建时间' })} value={formatTimestamp(record.created_at)} />
        </Box>
      </DetailSectionCard>

      <Box sx={{ mt: 3 }}>
        <DetailSectionCard
          title={translate('resources.voucher-batches.sections.vouchers.title', { _: '券码列表' })}
          icon={<VoucherIcon />}
        >
          <BatchVoucherList batchId={record.id} />
        </DetailSectionCard>
      </Box>
    </Box>
  );
};

export const VoucherBatchShow = () => {
  return (
    <Show>
      <VoucherBatchDetails />
    </Show>
  );
};
