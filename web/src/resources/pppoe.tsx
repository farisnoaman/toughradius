// web/src/resources/pppoe.tsx
import {
  List,
  Datagrid,
  TextField,
  DateField,
  Edit,
  TextInput,
  Create,
  BooleanInput,
  NumberInput,
  required,
  useRecordContext,
  Toolbar,
  SaveButton,
  DeleteButton,
  SimpleForm,
  TopToolbar,
  CreateButton,
  useTranslate,
  useListContext,
  FunctionField,
  SelectInput,
  ReferenceInput,
} from 'react-admin';
import {
  Box,
  Chip,
  Typography,
  Card,
  Avatar,
} from '@mui/material';
import { Router as PppoeIcon, CheckCircle as EnabledIcon, Cancel as DisabledIcon } from '@mui/icons-material';
import {
  ServerPagination,
  FormSection,
  FieldGrid,
  FieldGridItem,
  formLayoutSx,
  controlWrapperSx,
} from '../components';

const LARGE_LIST_PER_PAGE = 50;

// ============ 类型定义 ============

interface PppoeProfile {
  id: number;
  name?: string;
  node_id?: number;
  status?: 'enabled' | 'disabled';
  addr_pool?: string;
  ipv6_prefix_pool?: string;
  ipv6_addr_pool?: string;
  ac_name?: string;
  service_name?: string;
  session_timeout?: number;
  idle_timeout?: number;
  interim_interval?: number;
  up_rate?: number;
  down_rate?: number;
  up_burst_rate?: number;
  down_burst_rate?: number;
  up_burst_size?: number;
  down_burst_size?: number;
  vlanid1?: number;
  vlanid2?: number;
  pvc_vpi?: number;
  pvc_vci?: number;
  domain?: string;
  bind_mac?: number;
  bind_vlan?: number;
  active_num?: number;
  priority?: number;
  remark?: string;
  created_at?: string;
}

interface PppoeUser {
  id: number;
  node_id?: number;
  profile_id?: number;
  username?: string;
  password?: string;
  realname?: string;
  mobile?: string;
  email?: string;
  address?: string;
  mac_addr?: string;
  ip_addr?: string;
  ipv6_addr?: string;
  delegated_ipv6_prefix?: string;
  vlanid1?: number;
  vlanid2?: number;
  domain?: string;
  status?: 'enabled' | 'disabled' | 'expired';
  expire_time?: string;
  online_count?: number;
  total_session_time?: number;
  total_input_bytes?: number;
  total_output_bytes?: number;
  created_at?: string;
}

// ============ 工具函数 ============

const formatRate = (rate?: number): string => {
  if (!rate || rate === 0) return '-';
  if (rate >= 1024) {
    return `${(rate / 1024).toFixed(1)} Mbps`;
  }
  return `${rate} Kbps`;
};

// ============ 状态组件 ============

const StatusIndicator = ({ status }: { status?: string }) => {
  const translate = useTranslate();
  const isEnabled = status === 'enabled';
  return (
    <Chip
      icon={isEnabled ? <EnabledIcon sx={{ fontSize: '0.85rem !important' }} /> : <DisabledIcon sx={{ fontSize: '0.85rem !important' }} />}
      label={isEnabled ? translate('common.enabled', { _: '启用' }) : translate('common.disabled', { _: '禁用' })}
      size="small"
      color={isEnabled ? 'success' : 'default'}
      variant={isEnabled ? 'filled' : 'outlined'}
      sx={{ height: 22, fontWeight: 500, fontSize: '0.75rem' }}
    />
  );
};

// ============ 字段组件 ============

const ProfileNameField = () => {
  const record = useRecordContext<PppoeProfile>();
  if (!record) return null;

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <Avatar
        sx={{
          width: 32,
          height: 32,
          fontSize: '0.85rem',
          fontWeight: 600,
          bgcolor: record.status === 'enabled' ? 'info.main' : 'grey.400',
        }}
      >
        <PppoeIcon sx={{ fontSize: 18 }} />
      </Avatar>
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 600, color: 'text.primary', lineHeight: 1.3 }}>
          {record.name || '-'}
        </Typography>
        <StatusIndicator status={record.status} />
      </Box>
    </Box>
  );
};

const RateField = ({ source }: { source: 'up_rate' | 'down_rate' }) => {
  const record = useRecordContext<PppoeProfile>();
  if (!record) return null;

  return (
    <Chip
      label={formatRate(record[source])}
      size="small"
      color="info"
      variant="outlined"
      sx={{ fontFamily: 'monospace', fontSize: '0.8rem', height: 24 }}
    />
  );
};

const VlanField = () => {
  const record = useRecordContext<PppoeProfile>();
  if (!record) return null;

  const vlans = [record.vlanid1, record.vlanid2].filter(Boolean).join('/');
  return vlans ? (
    <Chip
      label={vlans}
      size="small"
      color="secondary"
      variant="outlined"
      sx={{ fontFamily: 'monospace', fontSize: '0.8rem', height: 24 }}
    />
  ) : <Typography variant="body2" color="text.secondary">-</Typography>;
};

const UsernameField = () => {
  const record = useRecordContext<PppoeUser>();
  if (!record) return null;

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <Avatar
        sx={{
          width: 32,
          height: 32,
          fontSize: '0.85rem',
          fontWeight: 600,
          bgcolor: record.status === 'enabled' ? 'success.main' : 'grey.400',
        }}
      >
        {record.username?.charAt(0).toUpperCase() || 'U'}
      </Avatar>
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 600, color: 'text.primary', lineHeight: 1.3 }}>
          {record.username || '-'}
        </Typography>
        <StatusIndicator status={record.status} />
      </Box>
    </Box>
  );
};

// ============ 列表操作栏 ============

const ProfileListActions = () => {
  const translate = useTranslate();
  return (
    <TopToolbar>
      <CreateButton label={translate('resources.pppoe-profiles.actions.create', { _: '新建策略' })} />
    </TopToolbar>
  );
};

const UserListActions = () => {
  const translate = useTranslate();
  return (
    <TopToolbar>
      <CreateButton label={translate('resources.pppoe-users.actions.create', { _: '新建用户' })} />
    </TopToolbar>
  );
};

// ============ 列表内容 ============

const ProfileListContent = () => {
  const translate = useTranslate();
  const { data, isLoading, total } = useListContext();

  if (isLoading) {
    return <Typography>Loading...</Typography>;
  }

  if (!data || data.length === 0) {
    return (
      <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}` }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 8, color: 'text.secondary' }}>
          <PppoeIcon sx={{ fontSize: 64, opacity: 0.3, mb: 2 }} />
          <Typography variant="h6" sx={{ opacity: 0.6 }}>
            {translate('resources.pppoe-profiles.empty.title', { _: '暂无策略' })}
          </Typography>
        </Box>
      </Card>
    );
  }

  return (
    <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}`, overflow: 'hidden' }}>
      <Box sx={{ px: 2, py: 1, bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.02)' : 'rgba(0,0,0,0.01)' }}>
        <Typography variant="body2" color="text.secondary">
          共 <strong>{total?.toLocaleString() || 0}</strong> 个策略
        </Typography>
      </Box>
      <Box sx={{ overflowX: 'auto' }}>
        <Datagrid rowClick="show" bulkActionButtons={false}>
          <FunctionField
            source="name"
            label={translate('resources.pppoe-profiles.fields.name', { _: '策略名称' })}
            render={() => <ProfileNameField />}
          />
          <TextField
            source="addr_pool"
            label={translate('resources.pppoe-profiles.fields.addr_pool', { _: '地址池' })}
          />
          <FunctionField
            source="up_rate"
            label={translate('resources.pppoe-profiles.fields.up_rate', { _: '上行速率' })}
            render={() => <RateField source="up_rate" />}
          />
          <FunctionField
            source="down_rate"
            label={translate('resources.pppoe-profiles.fields.down_rate', { _: '下行速率' })}
            render={() => <RateField source="down_rate" />}
          />
          <FunctionField
            source="vlanid1"
            label={translate('resources.pppoe-profiles.fields.vlan', { _: 'VLAN' })}
            render={() => <VlanField />}
          />
          <TextField
            source="active_num"
            label={translate('resources.pppoe-profiles.fields.active_num', { _: '并发数' })}
          />
          <DateField
            source="created_at"
            label={translate('resources.pppoe-profiles.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Box>
    </Card>
  );
};

// ============ 列表页面 ============

export const PppoeProfileList = () => {
  return (
    <List
      actions={<ProfileListActions />}
      sort={{ field: 'created_at', order: 'DESC' }}
      perPage={LARGE_LIST_PER_PAGE}
      pagination={<ServerPagination />}
      empty={false}
    >
      <ProfileListContent />
    </List>
  );
};

export const PppoeUserList = () => {
  const translate = useTranslate();
  return (
    <List
      actions={<UserListActions />}
      sort={{ field: 'created_at', order: 'DESC' }}
      perPage={LARGE_LIST_PER_PAGE}
      pagination={<ServerPagination />}
      empty={false}
    >
      <Card elevation={0} sx={{ borderRadius: 2, border: theme => `1px solid ${theme.palette.divider}`, overflow: 'hidden' }}>
        <Datagrid bulkActionButtons={false}>
          <FunctionField
            source="username"
            label={translate('resources.pppoe-users.fields.username', { _: '用户名' })}
            render={() => <UsernameField />}
          />
          <TextField
            source="realname"
            label={translate('resources.pppoe-users.fields.realname', { _: '姓名' })}
          />
          <TextField
            source="ip_addr"
            label={translate('resources.pppoe-users.fields.ip_addr', { _: 'IP地址' })}
          />
          <TextField
            source="mac_addr"
            label={translate('resources.pppoe-users.fields.mac_addr', { _: 'MAC地址' })}
          />
          <FunctionField
            source="expire_time"
            label={translate('resources.pppoe-users.fields.expire_time', { _: '过期时间' })}
            render={(record: PppoeUser) => record.expire_time ? new Date(record.expire_time).toLocaleDateString() : '-'}
          />
          <DateField
            source="created_at"
            label={translate('resources.pppoe-users.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Card>
    </List>
  );
};

// ============ 创建/编辑表单 ============

const ProfileFormToolbar = (props: any) => (
  <Toolbar {...props}>
    <SaveButton />
    <DeleteButton mutationMode="pessimistic" />
  </Toolbar>
);

export const PppoeProfileCreate = () => {
  const translate = useTranslate();

  return (
    <Create>
      <SimpleForm sx={formLayoutSx}>
        <FormSection
          title={translate('resources.pppoe-profiles.sections.basic.title', { _: '基本信息' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="name"
                label={translate('resources.pppoe-profiles.fields.name', { _: '策略名称' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput
                  source="status"
                  label={translate('resources.pppoe-profiles.fields.status_enabled', { _: '启用状态' })}
                  defaultValue={true}
                />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.pppoe-profiles.sections.network.title', { _: '网络设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 3 }}>
            <FieldGridItem>
              <TextInput
                source="addr_pool"
                label={translate('resources.pppoe-profiles.fields.addr_pool', { _: 'IPv4地址池' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="ipv6_prefix_pool"
                label={translate('resources.pppoe-profiles.fields.ipv6_prefix_pool', { _: 'IPv6前缀池' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="ipv6_addr_pool"
                label={translate('resources.pppoe-profiles.fields.ipv6_addr_pool', { _: 'IPv6地址池' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="ac_name"
                label={translate('resources.pppoe-profiles.fields.ac_name', { _: 'AC名称' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="service_name"
                label={translate('resources.pppoe-profiles.fields.service_name', { _: '服务名称' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="domain"
                label={translate('resources.pppoe-profiles.fields.domain', { _: '域名' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.pppoe-profiles.sections.time.title', { _: '时间设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 3 }}>
            <FieldGridItem>
              <NumberInput
                source="session_timeout"
                label={translate('resources.pppoe-profiles.fields.session_timeout', { _: '会话超时(秒)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="idle_timeout"
                label={translate('resources.pppoe-profiles.fields.idle_timeout', { _: '空闲超时(秒)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="interim_interval"
                label={translate('resources.pppoe-profiles.fields.interim_interval', { _: '计费间隔(秒)' })}
                min={0}
                defaultValue={600}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.pppoe-profiles.sections.bandwidth.title', { _: '带宽限制' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 4 }}>
            <FieldGridItem>
              <NumberInput
                source="up_rate"
                label={translate('resources.pppoe-profiles.fields.up_rate', { _: '上行速率(Kbps)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="down_rate"
                label={translate('resources.pppoe-profiles.fields.down_rate', { _: '下行速率(Kbps)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="up_burst_rate"
                label={translate('resources.pppoe-profiles.fields.up_burst_rate', { _: '上行突发速率' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="down_burst_rate"
                label={translate('resources.pppoe-profiles.fields.down_burst_rate', { _: '下行突发速率' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.pppoe-profiles.sections.vlan.title', { _: 'VLAN/PVC设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 4 }}>
            <FieldGridItem>
              <NumberInput
                source="vlanid1"
                label={translate('resources.pppoe-profiles.fields.vlanid1', { _: '内层VLAN' })}
                min={0}
                max={4096}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="vlanid2"
                label={translate('resources.pppoe-profiles.fields.vlanid2', { _: '外层VLAN' })}
                min={0}
                max={4096}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="pvc_vpi"
                label={translate('resources.pppoe-profiles.fields.pvc_vpi', { _: 'PVC VPI' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="pvc_vci"
                label={translate('resources.pppoe-profiles.fields.pvc_vci', { _: 'PVC VCI' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.pppoe-profiles.sections.session.title', { _: '会话设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 3 }}>
            <FieldGridItem>
              <NumberInput
                source="active_num"
                label={translate('resources.pppoe-profiles.fields.active_num', { _: '并发数' })}
                min={0}
                max={100}
                defaultValue={1}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="priority"
                label={translate('resources.pppoe-profiles.fields.priority', { _: 'QoS优先级' })}
                min={0}
                max={7}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput
                  source="bind_mac"
                  label={translate('resources.pppoe-profiles.fields.bind_mac', { _: '绑定MAC' })}
                  defaultValue={false}
                />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Create>
  );
};

export const PppoeProfileEdit = () => {
  return (
    <Edit>
      <SimpleForm toolbar={<ProfileFormToolbar />} sx={formLayoutSx}>
        <FormSection title="基本信息">
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput source="id" disabled fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput source="name" validate={[required()]} fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput source="status" label="启用状态" />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
        <FormSection title="网络设置">
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput source="addr_pool" fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput source="ipv6_prefix_pool" fullWidth size="small" />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
        <FormSection title="带宽限制">
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <NumberInput source="up_rate" min={0} fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput source="down_rate" min={0} fullWidth size="small" />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Edit>
  );
};

export const PppoeUserCreate = () => {
  const translate = useTranslate();

  return (
    <Create>
      <SimpleForm sx={formLayoutSx}>
        <FormSection title={translate('resources.pppoe-users.sections.basic.title', { _: '基本信息' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <ReferenceInput
                source="profile_id"
                reference="pppoe-profiles"
                label={translate('resources.pppoe-users.fields.profile_id', { _: 'PPPoE策略' })}
              >
                <SelectInput optionText="name" validate={[required()]} fullWidth size="small" />
              </ReferenceInput>
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="username"
                label={translate('resources.pppoe-users.fields.username', { _: '用户名' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="password"
                label={translate('resources.pppoe-users.fields.password', { _: '密码' })}
                type="password"
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="realname"
                label={translate('resources.pppoe-users.fields.realname', { _: '姓名' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.pppoe-users.sections.contact.title', { _: '联系方式' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="mobile"
                label={translate('resources.pppoe-users.fields.mobile', { _: '手机号' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="email"
                label={translate('resources.pppoe-users.fields.email', { _: '邮箱' })}
                type="email"
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem span={{ xs: 1, sm: 2 }}>
              <TextInput
                source="address"
                label={translate('resources.pppoe-users.fields.address', { _: '地址' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.pppoe-users.sections.network.title', { _: '网络设置' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="ip_addr"
                label={translate('resources.pppoe-users.fields.ip_addr', { _: '静态IPv4' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="ipv6_addr"
                label={translate('resources.pppoe-users.fields.ipv6_addr', { _: '静态IPv6' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="delegated_ipv6_prefix"
                label={translate('resources.pppoe-users.fields.delegated_ipv6_prefix', { _: 'IPv6前缀委派' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="mac_addr"
                label={translate('resources.pppoe-users.fields.mac_addr', { _: 'MAC地址' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="vlanid1"
                label={translate('resources.pppoe-users.fields.vlanid1', { _: '内层VLAN' })}
                min={0}
                max={4096}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="vlanid2"
                label={translate('resources.pppoe-users.fields.vlanid2', { _: '外层VLAN' })}
                min={0}
                max={4096}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.pppoe-users.sections.expiry.title', { _: '有效期' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="expire_time"
                label={translate('resources.pppoe-users.fields.expire_time', { _: '过期时间' })}
                type="datetime-local"
                fullWidth
                size="small"
                InputLabelProps={{ shrink: true }}
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Create>
  );
};

export const PppoeUserEdit = () => {
  return (
    <Edit>
      <SimpleForm toolbar={<ProfileFormToolbar />} sx={formLayoutSx}>
        <FormSection title="基本信息">
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput source="id" disabled fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput source="username" validate={[required()]} fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput source="password" type="password" fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput source="status" label="启用状态" />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
        <FormSection title="网络设置">
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput source="ip_addr" fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput source="mac_addr" fullWidth size="small" />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>
      </SimpleForm>
    </Edit>
  );
};
