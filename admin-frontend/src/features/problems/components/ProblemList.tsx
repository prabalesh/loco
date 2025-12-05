import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import type { ColumnsType } from "antd/es/table"
import { Button, Table, Tag, Popconfirm } from "antd"
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons"
import dayjs from "dayjs"
import toast from "react-hot-toast"
import { useNavigate } from "react-router-dom"
import type { Problem } from "../../../types"
import { adminProblemApi } from "../../../api/adminApi"

export default function ProblemList() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data, isFetching } = useQuery({
    queryKey: ["admin-problems"],
    queryFn: async () => {
      const res = await adminProblemApi.getAll()
      return res.data
    },
  })

  const problems = data?.data || []

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminProblemApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] })
      toast.success("Problem deleted successfully")
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to delete problem")
    },
  })

  const columns: ColumnsType<Problem> = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 70,
      sorter: (a, b) => a.id - b.id,
      render: (id) => <span className="font-mono text-gray-600">#{id}</span>,
    },
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
      render: (title, record) => (
        <div>
          <div className="font-semibold text-gray-900">{title}</div>
          <div className="text-xs text-gray-500 font-mono">{record.slug}</div>
        </div>
      ),
    },
    {
      title: "Difficulty",
      dataIndex: "difficulty",
      key: "difficulty",
      width: 120,
      filters: [
        { text: "Easy", value: "easy" },
        { text: "Medium", value: "medium" },
        { text: "Hard", value: "hard" },
      ],
      onFilter: (value, record) => record.difficulty === value,
      render: (difficulty) => {
        const color =
          difficulty === "easy" ? "green" : difficulty === "medium" ? "orange" : "red"
        return <Tag color={color}>{difficulty.toUpperCase()}</Tag>
      },
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: 120,
      filters: [
        { text: "Draft", value: "draft" },
        { text: "Published", value: "published" },
      ],
      onFilter: (value, record) => record.status === value,
      render: (status) => (
        <Tag color={status === "published" ? "blue" : "default"}>
          {status.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: "Active",
      dataIndex: "is_active",
      key: "is_active",
      width: 100,
      render: (isActive: boolean) => (
        <Tag color={isActive ? "green" : "red"}>
          {isActive ? "YES" : "NO"}
        </Tag>
      ),
    },
    {
      title: "Limits",
      key: "limits",
      width: 160,
      render: (_, record) => (
        <div className="text-xs text-gray-600">
          <div>⏱️ Time: {record.time_limit}ms</div>
          <div>💾 Memory: {record.memory_limit}MB</div>
        </div>
      ),
    },
    {
      title: "Acceptance",
      key: "acceptance",
      width: 160,
      sorter: (a, b) => a.acceptance_rate - b.acceptance_rate,
      render: (_, record) => (
        <div className="text-xs text-gray-600">
          <div className="font-semibold">
            {record.acceptance_rate.toFixed(1)}%
          </div>
          <div>
            {record.total_accepted}/{record.total_submissions} AC
          </div>
        </div>
      ),
    },
    {
      title: "Created",
      dataIndex: "created_at",
      key: "created_at",
      width: 150,
      sorter: (a, b) => dayjs(a.created_at).unix() - dayjs(b.created_at).unix(),
      render: (date) => (
        <span className="text-xs text-gray-600">
          {dayjs(date).format("MMM DD, YYYY")}
        </span>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      width: 150,
      fixed: "right",
      render: (_, record) => (
        <div className="flex gap-2">
          <Button
            type="primary"
            ghost
            size="small"
            icon={<EditOutlined />}
            onClick={() => navigate(`/problems/edit/${record.id}`)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete Problem"
            description={`Delete "${record.title}"?`}
            icon={<ExclamationCircleOutlined style={{ color: "red" }} />}
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="Delete"
            okButtonProps={{ danger: true }}
            cancelText="Cancel"
          >
            <Button
              danger
              size="small"
              icon={<DeleteOutlined />}
              loading={deleteMutation.isPending}
            />
          </Popconfirm>
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="flex justify-between items-center mb-6 p-5 bg-white rounded-lg shadow">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Problem Management</h1>
          <p className="text-sm text-gray-500 mt-1">
            Manage coding problems for your platform
          </p>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate("/problems/create")}
          size="large"
        >
          Create Problem
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <Table
          bordered
          rowKey="id"
          columns={columns}
          dataSource={problems}
          loading={isFetching || deleteMutation.isPending}
          pagination={{
            pageSize: data?.limit || 10,
            total: data?.total,
            showSizeChanger: true,
            showTotal: (total) => `Total ${total} problems`,
          }}
          scroll={{ x: 1400 }}
          size="middle"
        />
      </div>
    </div>
  )
}
