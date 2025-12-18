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
import { PROBLEM_STEPS, ROUTES } from "../../../config/constant"
import { Play } from "lucide-react"

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
      title: "Title",
      dataIndex: "title",
      key: "title",
      width: 150,
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
      title: "Current Step",
      dataIndex: "current_step",
      key: "current_step",
      width: 120,
      filters: [
        { text: "Metadata", value: 1 },
        { text: "Testcases", value: 2 },
        { text: "Languages", value: 3 },
        { text: "Validate", value: 4 },
      ],
      onFilter: (value, record) => record.status === value,
      render: (status) => (
        <Tag>
          {PROBLEM_STEPS[status-1].label}
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
      title: "Updated",
      dataIndex: "updated_at",
      key: "updated_at",
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
          <Button
            type="primary"
            ghost
            size="small"
            icon={<Play className="h-4 w-4" />}
            onClick={() => {
              let link = "";
              switch(record.current_step) {
                case 1:
                  link = ROUTES.PROBLEMS.TESTCASES(record.id)
                  break
                case 2:
                  link = ROUTES.PROBLEMS.LANGUAGES(record.id)
                  break
                case 3:
                  link = ROUTES.PROBLEMS.VALIDATE(record.id)
                  break
                case 4:
                  link = ROUTES.PROBLEMS.VALIDATE(record.id)
                  break
              }
              navigate(link)
            }}
          >
            Resume
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
