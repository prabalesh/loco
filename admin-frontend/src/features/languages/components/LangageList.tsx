import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { adminLanguagesApi } from "../../../api/adminApi"
import type { ColumnsType } from "antd/es/table"
import type { Language } from "../../../types"
import Table from "antd/es/table"
import { CheckCircleOutlined, CloseCircleOutlined } from "@ant-design/icons"
import { Button, Switch } from "antd"
import dayjs from "dayjs"
import toast from "react-hot-toast"

export default function LanguageList() {
    const queryClient = useQueryClient()
    const {data, isFetching, refetch } = useQuery({
        queryKey: ["admin-languages"],
        queryFn: async () => {
            const response = await adminLanguagesApi.getAll()
            return response.data
        }
    })
    const languages = data?.data || []

    const updateStatusMutation = useMutation({
        mutationFn: ({languageId, activate}: {languageId: number; activate: Boolean}) => {
            return activate?adminLanguagesApi.deactivate(languageId):adminLanguagesApi.activate(languageId);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language status updated successfully")
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to update status')
        },
    })

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
            render: (date) => <span className="font-mono text-gray-600">{dayjs(date).format('MMM DD, YYYY HH:mm')}</span>,
        },
    ]

    return (
    <div>
        <div className="clas">
            <div>
                <h1 className="text-3xl font-bold">Language Management</h1>
            </div>
        </div>
        <div>
            <Button
                onClick={() => refetch()}
                loading={isFetching}
            >
                Refressh
            </Button>
            <Table 
                bordered
                rowKey="id"
                columns={columns}
                dataSource={languages}
                loading={isFetching}
                pagination={{ pageSize: 10, showSizeChanger: true }}
                scroll={{ x: 900 }}
                className="shadow-lg rounded-lg"
                size="middle"
            />
        </div>
    </div>)
}