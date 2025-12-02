import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table,
  Button,
  Space,
  Tag,
  Popconfirm,
  Select,
  Switch,
  Tooltip,
} from 'antd'
import {
  DeleteOutlined,
  EditOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons'
import { adminUsersApi } from '../../../api/adminApi'
import toast from 'react-hot-toast'
import dayjs from 'dayjs'
import type { User } from '../../../types'
import type { ColumnsType } from 'antd/es/table'

export const UsersList = () => {
  const queryClient = useQueryClient()
  const [editingRole, setEditingRole] = useState<{ userId: number; role: string } | null>(null)

  const { data: users, isLoading } = useQuery({
    queryKey: ['admin-users'],
    queryFn: async () => {
      const response = await adminUsersApi.getAll()
      return response.data
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (userId: number) => adminUsersApi.deleteUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('User deleted successfully')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to delete user')
    },
  })

  const updateRoleMutation = useMutation({
    mutationFn: ({ userId, role }: { userId: number; role: string }) =>
      adminUsersApi.updateRole(userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('Role updated successfully')
      setEditingRole(null)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update role')
    },
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ userId, isActive }: { userId: number; isActive: boolean }) =>
      adminUsersApi.updateStatus(userId, isActive),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('Status updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update status')
    },
  })

  const roleColors: Record<string, string> = {
    admin: 'red',
    moderator: 'orange',
    user: 'blue',
  }

  const columns: ColumnsType<User> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
      sorter: (a, b) => a.id - b.id,
      render: (id) => <span className="font-mono text-gray-600">{id}</span>,
    },
    {
      title: 'Username',
      dataIndex: 'username',
      key: 'username',
      sorter: (a, b) => a.username.localeCompare(b.username),
      render: (username) => <span className="font-semibold text-gray-800">{username}</span>,
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
      render: (email) => (
        <Tooltip title={email}>
          <span className="truncate block max-w-xs text-gray-700">{email}</span>
        </Tooltip>
      ),
    },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
      width: 180,
      render: (role: string, record: User) => {
        const isEditing = editingRole?.userId === record.id

        if (isEditing) {
          return (
            <Space>
              <Select
                value={editingRole.role}
                onChange={(value) => setEditingRole({ userId: record.id, role: value })}
                style={{ width: 140 }}
                options={[
                  { value: 'user', label: 'User' },
                  { value: 'admin', label: 'Admin' },
                  { value: 'moderator', label: 'Mod' },
                ]}
                disabled={updateRoleMutation.isPending}
              />
              <Button
                size="small"
                type="primary"
                onClick={() => updateRoleMutation.mutate({ userId: record.id, role: editingRole.role })}
                loading={updateRoleMutation.isPending}
                disabled={updateRoleMutation.isPending}
              >
                Save
              </Button>
              <Button size="small" onClick={() => setEditingRole(null)} disabled={updateRoleMutation.isPending}>
                Cancel
              </Button>
            </Space>
          )
        }

        return (
          <Space>
            <Tag color={roleColors[role] || 'default'} className="uppercase font-semibold tracking-wide">
              {role}
            </Tag>
            <Button
              size="small"
              icon={<EditOutlined />}
              onClick={() => setEditingRole({ userId: record.id, role })}
            />
          </Space>
        )
      },
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      width: 120,
      render: (isActive: boolean, record: User) => (
        <Switch
          checked={isActive}
          checkedChildren={<CheckCircleOutlined />}
          unCheckedChildren={<CloseCircleOutlined />}
          onChange={(checked) =>
            updateStatusMutation.mutate({ userId: record.id, isActive: checked })
          }
          loading={updateStatusMutation.isPending}
        />
      ),
    },
    {
      title: 'Email Verified',
      dataIndex: 'email_verified',
      key: 'email_verified',
      width: 140,
      render: (verified: boolean) =>
        verified ? (
          <Tag color="success" className="uppercase font-semibold">
            Verified
          </Tag>
        ) : (
          <Tag color="warning" className="uppercase font-semibold">
            Not Verified
          </Tag>
        ),
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      sorter: (a, b) => dayjs(a.created_at).unix() - dayjs(b.created_at).unix(),
      render: (date: string) => (
        <span className="text-gray-600 font-mono">{dayjs(date).format('MMM DD, YYYY HH:mm')}</span>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, record: User) => (
        <Popconfirm
          title="Delete user"
          description="Are you sure you want to delete this user?"
          onConfirm={() => deleteMutation.mutate(record.id)}
          okText="Yes"
          cancelText="No"
          disabled={deleteMutation.isPending}
        >
          <Button
            danger
            icon={<DeleteOutlined />}
            loading={deleteMutation.isPending}
            block
          >
            Delete
          </Button>
        </Popconfirm>
      ),
    },
  ]

  return (
    <div className="max-w-full">
      <h1 className="text-3xl font-semibold text-gray-900 mb-8">User Management</h1>
      <Button
        type="default"
        onClick={() => queryClient.invalidateQueries({ queryKey: ['admin-users'] })}
        loading={isLoading}
      >
        Refresh
      </Button>
      <Table
        bordered
        rowKey="id"
        columns={columns}
        dataSource={users}
        loading={isLoading}
        pagination={{ pageSize: 10, showSizeChanger: true }}
        scroll={{ x: 900 }}
        className="shadow-lg rounded-lg"
        size="middle"
      />
    </div>
  )
}
