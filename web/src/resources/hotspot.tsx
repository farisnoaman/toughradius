// web/src/resources/hotspot.tsx
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
import {
  Wifi as HotspotIcon,
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
} from '../components';

const LARGE_LIST_PER_PAGE = 50;

// ============ 类型定义 ============

interface HotspotProfile {
  id: number;
  name?: string;
  node_id?: number;
  status?: 'enabled' | 'disabled';
  auth_mode?: 'userpass' | 'mac' | 'mac-userpass';
  session_timeout?: number;
  idle_timeout?: number;
  daily_limit?: number;
  monthly_limit?: number;
  up_rate?: number;
  down_rate?: number;
  up_limit?: number;
  down_limit?: number;
  total_limit?: number;
  addr_pool?: string;
  domain?: string;
  welcome_url?: string;
  logout_url?: string;
  bind_mac?: number;
  max_devices?: number;
  remark?: string;
  created_at?: string;
}

interface HotspotUser {
  id: number;
  node_id?: number;
  profile_id?: number;
  username?: string;
  password?: string;
  realname?: string;
  mobile?: string;
  email?: string;
  mac_addr?: string;
  ip_addr?: string;
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

const AuthModeIndicator = ({ mode }: { mode?: string }) => {
  const translate = useTranslate();
  const modeConfig: Record<string, { label: string; color: 'primary' | 'secondary' | 'info' }> = {
    userpass: { label: translate('resources.hotspot-profiles.auth_mode.userpass', { _: '用户密码' }), color: 'primary' },
    mac: { label: translate('resources.hotspot-profiles.auth_mode.mac', { _: 'MAC认证' }), color: 'secondary' },
    'mac-userpass': { label: translate('resources.hotspot-profiles.auth_mode.mac-userpass', { _: 'MAC+密码' }), color: 'info' },
  };
  const config = modeConfig[mode || 'userpass'];
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

const ProfileNameField = () => {
  const record = useRecordContext<HotspotProfile>();
  if (!record) return null;

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <Avatar
        sx={{
          width: 32,
          height: 32,
          fontSize: '0.85rem',
          fontWeight: 600,
          bgcolor: record.status === 'enabled' ? 'secondary.main' : 'grey.400',
        }}
      >
        <HotspotIcon sx={{ fontSize: 18 }} />
      </Avatar>
      <Box>
        <Typography variant="body2" sx={{ fontWeight: 600, color: 'text.primary', lineHeight: 1.3 }}>
          {record.name || '-'}
        </Typography>
        <AuthModeIndicator mode={record.auth_mode} />
      </Box>
    </Box>
  );
};

const RateField = ({ source }: { source: 'up_rate' | 'down_rate' }) => {
  const record = useRecordContext<HotspotProfile>();
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

const UsernameField = () => {
  const record = useRecordContext<HotspotUser>();
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
      <CreateButton label={translate('resources.hotspot-profiles.actions.create', { _: '新建策略' })} />
    </TopToolbar>
  );
};

const UserListActions = () => {
  const translate = useTranslate();
  return (
    <TopToolbar>
      <CreateButton label={translate('resources.hotspot-users.actions.create', { _: '新建用户' })} />
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
          <HotspotIcon sx={{ fontSize: 64, opacity: 0.3, mb: 2 }} />
          <Typography variant="h6" sx={{ opacity: 0.6 }}>
            {translate('resources.hotspot-profiles.empty.title', { _: '暂无策略' })}
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
            label={translate('resources.hotspot-profiles.fields.name', { _: '策略名称' })}
            render={() => <ProfileNameField />}
          />
          <FunctionField
            source="up_rate"
            label={translate('resources.hotspot-profiles.fields.up_rate', { _: '上行速率' })}
            render={() => <RateField source="up_rate" />}
          />
          <FunctionField
            source="down_rate"
            label={translate('resources.hotspot-profiles.fields.down_rate', { _: '下行速率' })}
            render={() => <RateField source="down_rate" />}
          />
          <TextField
            source="session_timeout"
            label={translate('resources.hotspot-profiles.fields.session_timeout', { _: '会话时长(分)' })}
          />
          <TextField
            source="max_devices"
            label={translate('resources.hotspot-profiles.fields.max_devices', { _: '最大设备数' })}
          />
          <DateField
            source="created_at"
            label={translate('resources.hotspot-profiles.fields.created_at', { _: '创建时间' })}
            showTime
          />
        </Datagrid>
      </Box>
    </Card>
  );
};

// ============ 列表页面 ============

export const HotspotProfileList = () => {
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

export const HotspotUserList = () => {
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
            label={translate('resources.hotspot-users.fields.username', { _: '用户名' })}
            render={() => <UsernameField />}
          />
          <TextField
            source="realname"
            label={translate('resources.hotspot-users.fields.realname', { _: '姓名' })}
          />
          <TextField
            source="mac_addr"
            label={translate('resources.hotspot-users.fields.mac_addr', { _: 'MAC地址' })}
          />
          <FunctionField
            source="expire_time"
            label={translate('resources.hotspot-users.fields.expire_time', { _: '过期时间' })}
            render={(record: HotspotUser) => record.expire_time ? new Date(record.expire_time).toLocaleDateString() : '-'}
          />
          <DateField
            source="created_at"
            label={translate('resources.hotspot-users.fields.created_at', { _: '创建时间' })}
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

const authModeChoices = [
  { id: 'userpass', name: '用户名密码认证' },
  { id: 'mac', name: 'MAC地址认证' },
  { id: 'mac-userpass', name: 'MAC+密码双重认证' },
];

export const HotspotProfileCreate = () => {
  const translate = useTranslate();

  return (
    <Create>
      <SimpleForm sx={formLayoutSx}>
        <FormSection
          title={translate('resources.hotspot-profiles.sections.basic.title', { _: '基本信息' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="name"
                label={translate('resources.hotspot-profiles.fields.name', { _: '策略名称' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <SelectInput
                source="auth_mode"
                label={translate('resources.hotspot-profiles.fields.auth_mode', { _: '认证模式' })}
                choices={authModeChoices}
                defaultValue="userpass"
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.hotspot-profiles.sections.time.title', { _: '时间限制' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2, md: 3 }}>
            <FieldGridItem>
              <NumberInput
                source="session_timeout"
                label={translate('resources.hotspot-profiles.fields.session_timeout', { _: '会话时长(分钟)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="idle_timeout"
                label={translate('resources.hotspot-profiles.fields.idle_timeout', { _: '空闲超时(分钟)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="daily_limit"
                label={translate('resources.hotspot-profiles.fields.daily_limit', { _: '每日限额(分钟)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.hotspot-profiles.sections.bandwidth.title', { _: '带宽限制' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <NumberInput
                source="up_rate"
                label={translate('resources.hotspot-profiles.fields.up_rate', { _: '上行速率(Kbps)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <NumberInput
                source="down_rate"
                label={translate('resources.hotspot-profiles.fields.down_rate', { _: '下行速率(Kbps)' })}
                min={0}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.hotspot-profiles.sections.device.title', { _: '设备设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <NumberInput
                source="max_devices"
                label={translate('resources.hotspot-profiles.fields.max_devices', { _: '最大设备数' })}
                min={0}
                max={100}
                defaultValue={1}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput
                  source="bind_mac"
                  label={translate('resources.hotspot-profiles.fields.bind_mac', { _: '绑定MAC' })}
                  defaultValue={false}
                />
              </Box>
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection
          title={translate('resources.hotspot-profiles.sections.network.title', { _: '网络设置' })}
        >
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="addr_pool"
                label={translate('resources.hotspot-profiles.fields.addr_pool', { _: '地址池' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="domain"
                label={translate('resources.hotspot-profiles.fields.domain', { _: '域名' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem span={{ xs: 1, sm: 2 }}>
              <TextInput
                source="welcome_url"
                label={translate('resources.hotspot-profiles.fields.welcome_url', { _: '欢迎页面URL' })}
                type="url"
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem span={{ xs: 1, sm: 2 }}>
              <TextInput
                source="logout_url"
                label={translate('resources.hotspot-profiles.fields.logout_url', { _: '登出页面URL' })}
                type="url"
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

export const HotspotProfileEdit = () => {
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
              <SelectInput source="auth_mode" choices={authModeChoices} fullWidth size="small" />
            </FieldGridItem>
            <FieldGridItem>
              <Box sx={controlWrapperSx}>
                <BooleanInput source="status" label="启用状态" />
              </Box>
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

export const HotspotUserCreate = () => {
  const translate = useTranslate();

  return (
    <Create>
      <SimpleForm sx={formLayoutSx}>
        <FormSection title={translate('resources.hotspot-users.sections.basic.title', { _: '基本信息' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <ReferenceInput
                source="profile_id"
                reference="hotspot-profiles"
                label={translate('resources.hotspot-users.fields.profile_id', { _: '热点策略' })}
              >
                <SelectInput optionText="name" validate={[required()]} fullWidth size="small" />
              </ReferenceInput>
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="username"
                label={translate('resources.hotspot-users.fields.username', { _: '用户名' })}
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="password"
                label={translate('resources.hotspot-users.fields.password', { _: '密码' })}
                type="password"
                validate={[required()]}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="realname"
                label={translate('resources.hotspot-users.fields.realname', { _: '姓名' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.hotspot-users.sections.contact.title', { _: '联系方式' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="mobile"
                label={translate('resources.hotspot-users.fields.mobile', { _: '手机号' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="email"
                label={translate('resources.hotspot-users.fields.email', { _: '邮箱' })}
                type="email"
                fullWidth
                size="small"
              />
            </FieldGridItem>
          </FieldGrid>
        </FormSection>

        <FormSection title={translate('resources.hotspot-users.sections.device.title', { _: '设备信息' })}>
          <FieldGrid columns={{ xs: 1, sm: 2 }}>
            <FieldGridItem>
              <TextInput
                source="mac_addr"
                label={translate('resources.hotspot-users.fields.mac_addr', { _: 'MAC地址' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="ip_addr"
                label={translate('resources.hotspot-users.fields.ip_addr', { _: 'IP地址' })}
                fullWidth
                size="small"
              />
            </FieldGridItem>
            <FieldGridItem>
              <TextInput
                source="expire_time"
                label={translate('resources.hotspot-users.fields.expire_time', { _: '过期时间' })}
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

export const HotspotUserEdit = () => {
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
      </SimpleForm>
    </Edit>
  );
};
