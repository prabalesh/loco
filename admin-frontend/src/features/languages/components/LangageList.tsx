import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { adminLanguagesApi } from "../../../api/adminApi"
import type { ColumnsType } from "antd/es/table"
import type { Language } from "../../../types"
import { Modal, Form, Input, Button, Table, Switch, Popconfirm } from "antd"
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckCircleOutlined, CloseCircleOutlined, ExclamationCircleOutlined } from "@ant-design/icons"
import dayjs from "dayjs"
import toast from "react-hot-toast"
import { useState } from "react"
import type { CreateOrUpdateLanguageRequest } from "../../../types/request"

interface LanguageFormValues extends CreateOrUpdateLanguageRequest {}

export default function LanguageList() {
    const queryClient = useQueryClient()
    const [form] = Form.useForm<LanguageFormValues>()
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [editingLanguage, setEditingLanguage] = useState<Language | null>(null)
    
    const { data, isFetching, refetch } = useQuery({
        queryKey: ["admin-languages"],
        queryFn: async () => {
            const response = await adminLanguagesApi.getAll()
            return response.data
        }
    })
    
    const languages = data?.data || []

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (values: CreateOrUpdateLanguageRequest) => adminLanguagesApi.create(values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language created successfully")
            setIsModalOpen(false)
            form.resetFields()
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to create language')
        },
    })

    // Update mutation
    const updateMutation = useMutation({
        mutationFn: ({ id, values }: { id: number; values: CreateOrUpdateLanguageRequest }) => 
            adminLanguagesApi.update(id, values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language updated successfully")
            setIsModalOpen(false)
            form.resetFields()
            setEditingLanguage(null)
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to update language')
        },
    })

    // Status update mutation
    const updateStatusMutation = useMutation({
        mutationFn: ({languageId, activate}: {languageId: number; activate: boolean}) => {
            return activate ? adminLanguagesApi.deactivate(languageId) : adminLanguagesApi.activate(languageId)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language status updated successfully")
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to update status')
        },
    })

    // DELETE mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => adminLanguagesApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language deleted successfully")
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to delete language')
        },
    })

    const handleOpenModal = (language?: Language) => {
        if (language) {
            form.setFieldsValue({
                language_id: language.language_id,
                name: language.name,
                version: language.version,
                extension: language.extension,
                default_template: language.default_template || "",
            })
            setEditingLanguage(language)
        } else {
            form.resetFields()
            setEditingLanguage(null)
        }
        setIsModalOpen(true)
    }

    const handleSubmit = (values: LanguageFormValues) => {
        if (editingLanguage) {
            updateMutation.mutate({ id: editingLanguage.id, values })
        } else {
            createMutation.mutate(values)
        }
    }

    const handleDelete = (id: number) => {
        deleteMutation.mutate(id)
    }

    const columns: ColumnsType<Language> = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 70,
            sorter: (a, b) => a.id - b.id,
            render: (id) => <span className="font-mono text-gray-600">{id}</span>,
        },
        {
            title: 'Language ID',
            dataIndex: 'language_id',
            key: 'language_id',
            render: (language_id) => <span className="font-mono text-gray-600">{language_id}</span>,
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
            render: (name) => <span className="font-mono text-gray-600">{name}</span>,
        },
        {
            title: 'Status',
            dataIndex: 'is_active',
            key: 'is_active',
            width: 120,
            render: (isActive: boolean, record: Language) => (
                <Switch
                    checked={isActive}
                    checkedChildren={<CheckCircleOutlined />}
                    unCheckedChildren={<CloseCircleOutlined />}
                    onChange={() => 
                        updateStatusMutation.mutate({languageId: record.id, activate: record.is_active})
                    }
                    loading={updateStatusMutation.isPending}
                />
            ),
        },
        {
            title: 'Extension',
            dataIndex: 'extension',
            key: 'extension',
            render: (extension) => <span className="font-mono text-gray-600">{extension}</span>,
        },
        {
            title: 'Version',
            dataIndex: 'version',
            key: 'version',
            render: (version) => <span className="font-mono text-gray-600">{version}</span>,
        },
        {
            title: 'Template',
            dataIndex: 'template',
            key: 'template',
            render: (template) => <span className="font-mono text-gray-600 truncate max-w-[150px]">{template}</span>,
        },
        {
            title: 'Created At',
            dataIndex: 'created_at',
            key: 'created_at',
            width: 180,
            sorter: (a, b) => dayjs(a.created_at).unix() - dayjs(b.created_at).unix(),
            render: (date) => <span className="font-mono text-gray-600">{dayjs(date).format('MMM DD, YYYY HH:mm')}</span>,
        },
        {
            title: 'Updated At',
            dataIndex: 'updated_at',
            key: 'updated_at',
            width: 180,
            render: (date) => <span className="font-mono text-gray-600">{dayjs(date).format('MMM DD, YYYY HH:mm')}</span>,
        },
        {
            title: 'Actions',
            key: 'actions',
            width: 160,
            render: (_, record: Language) => (
                <div className="flex space-x-1">
                    <Button 
                        type="link" 
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleOpenModal(record)}
                        disabled={updateStatusMutation.isPending || deleteMutation.isPending}
                    >
                        Edit
                    </Button>
                    <Popconfirm
                        title="Delete Language"
                        description={`Are you sure you want to delete "${record.name}"?`}
                        icon={<ExclamationCircleOutlined style={{ color: 'red' }} />}
                        onConfirm={() => handleDelete(record.id)}
                        okText="Yes, Delete"
                        okButtonProps={{ danger: true }}
                        cancelText="Cancel"
                        disabled={deleteMutation.isPending}
                    >
                        <Button 
                            type="link" 
                            size="small"
                            icon={<DeleteOutlined />}
                            loading={deleteMutation.isPending}
                            disabled={deleteMutation.isPending}
                        />
                    </Popconfirm>
                </div>
            ),
        },
    ]

    return (
        <div>
            <div className="flex justify-between items-center mb-6 p-4 bg-white rounded-lg shadow-sm">
                <h1 className="text-3xl font-bold">Language Management</h1>
                <Button 
                    type="primary" 
                    icon={<PlusOutlined />}
                    onClick={() => handleOpenModal()}
                    loading={createMutation.isPending}
                    disabled={deleteMutation.isPending}
                >
                    Add Language
                </Button>
            </div>
            
            <div className="bg-white rounded-lg shadow-lg">
                <Button
                    onClick={() => refetch()}
                    loading={isFetching}
                    className="mb-4 ml-4 mt-4"
                    disabled={deleteMutation.isPending}
                >
                    Refresh
                </Button>
                <Table 
                    bordered
                    rowKey="id"
                    columns={columns}
                    dataSource={languages}
                    loading={isFetching || createMutation.isPending || updateMutation.isPending || updateStatusMutation.isPending || deleteMutation.isPending}
                    pagination={{ pageSize: 10, showSizeChanger: true }}
                    scroll={{ x: 1600 }}
                    className="shadow-lg rounded-lg"
                    size="middle"
                />
            </div>

            {/* CREATE/UPDATE MODAL */}
            <Modal
                title={editingLanguage ? "Edit Language" : "Create Language"}
                open={isModalOpen}
                onCancel={() => {
                    setIsModalOpen(false)
                    form.resetFields()
                    setEditingLanguage(null)
                }}
                footer={null}
                width={600}
                destroyOnClose
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                    disabled={createMutation.isPending || updateMutation.isPending}
                >
                    <Form.Item
                        name="language_id"
                        label="Language ID"
                        rules={[{ required: true, message: 'Please enter language ID' }]}
                    >
                        <Input placeholder="e.g., en, fr, es" />
                    </Form.Item>
                    
                    <Form.Item
                        name="name"
                        label="Name"
                        rules={[{ required: true, message: 'Please enter language name' }]}
                    >
                        <Input placeholder="Enter language name" />
                    </Form.Item>
                    
                    <Form.Item
                        name="version"
                        label="Version"
                        rules={[{ required: true, message: 'Please enter version' }]}
                    >
                        <Input placeholder="e.g., 1.0.0, latest" />
                    </Form.Item>
                    
                    <Form.Item
                        name="extension"
                        label="File Extension"
                        rules={[{ required: true, message: 'Please enter file extension' }]}
                    >
                        <Input placeholder="e.g., .js, .py, .java" />
                    </Form.Item>
                    
                    <Form.Item
                        name="template"
                        label="Default Template"
                        rules={[{ required: true, message: 'Please enter default template' }]}
                    >
                        <Input.TextArea rows={4} placeholder="Enter default code template" />
                    </Form.Item>
                    
                    <Form.Item className="mb-0">
                        <div className="flex justify-end space-x-2">
                            <Button 
                                onClick={() => {
                                    setIsModalOpen(false)
                                    form.resetFields()
                                    setEditingLanguage(null)
                                }}
                            >
                                Cancel
                            </Button>
                            <Button 
                                type="primary" 
                                htmlType="submit"
                                loading={createMutation.isPending || updateMutation.isPending}
                            >
                                {editingLanguage ? 'Update Language' : 'Create Language'}
                            </Button>
                        </div>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}
