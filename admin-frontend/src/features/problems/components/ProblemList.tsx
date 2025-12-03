import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import type { ColumnsType } from "antd/es/table"
import type { CreateOrUpdateProblemRequest } from "../../../types/request"
import {
  Modal,
  Form,
  Input,
  Button,
  Table,
  Switch,
  Select,
  InputNumber,
  Tag,
  Popconfirm,
} from "antd"
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons"
import dayjs from "dayjs"
import toast from "react-hot-toast"
import { useState } from "react"
import type { Problem } from "../../../types"
import { adminProblemApi } from "../../../api/adminApi"

interface ProblemFormValues extends CreateOrUpdateProblemRequest {}

export default function ProblemList() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<ProblemFormValues>()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProblem, setEditingProblem] = useState<Problem | null>(null)

  const { data, isFetching } = useQuery({
    queryKey: ["admin-problems"],
    queryFn: async () => {
      const res = await adminProblemApi.getAll()
      return res.data // PaginatedResponse<Problem>
    },
  })

  const problems = data?.data || []

  const createMutation = useMutation({
    mutationFn: (values: CreateOrUpdateProblemRequest) =>
      adminProblemApi.create(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] })
      toast.success("Problem created")
      setIsModalOpen(false)
      form.resetFields()
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to create problem")
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: number; values: CreateOrUpdateProblemRequest }) =>
      adminProblemApi.update(id, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] })
      toast.success("Problem updated")
      setIsModalOpen(false)
      form.resetFields()
      setEditingProblem(null)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to update problem")
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminProblemApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] })
      toast.success("Problem deleted")
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to delete problem")
    },
  })

  const handleOpenModal = (problem?: Problem) => {
    if (problem) {
      form.setFieldsValue({
        title: problem.title,
        slug: problem.slug,
        description: problem.description,
        difficulty: problem.difficulty,
        time_limit: problem.time_limit,
        memory_limit: problem.memory_limit,
        validator_type: problem.validator_type,
        input_format: problem.input_format,
        output_format: problem.output_format,
        constraints: problem.constraints,
        status: problem.status,
        is_active: problem.is_active,
      })
      setEditingProblem(problem)
    } else {
      form.resetFields()
      form.setFieldsValue({
        difficulty: "easy",
        validator_type: "exact_match",
        status: "draft",
        is_active: true,
      })
      setEditingProblem(null)
    }
    setIsModalOpen(true)
  }

  const handleSubmit = (values: ProblemFormValues) => {
    if (editingProblem) {
      updateMutation.mutate({ id: editingProblem.id, values })
    } else {
      createMutation.mutate(values)
    }
  }

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id)
  }

  const columns: ColumnsType<Problem> = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 70,
      sorter: (a, b) => a.id - b.id,
      render: (id) => <span className="font-mono text-gray-600">{id}</span>,
    },
    {
      title: "Title",
      dataIndex: "title",
      key: "title",
      render: (title, record) => (
        <div>
          <div className="font-semibold">{title}</div>
          <div className="text-xs text-gray-500">{record.slug}</div>
        </div>
      ),
    },
    {
      title: "Difficulty",
      dataIndex: "difficulty",
      key: "difficulty",
      width: 120,
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
      width: 110,
      render: (isActive: boolean) => (
        <Tag color={isActive ? "green" : "red"}>
          {isActive ? "ACTIVE" : "INACTIVE"}
        </Tag>
      ),
    },
    {
      title: "Limits",
      key: "limits",
      width: 160,
      render: (_, record) => (
        <div className="text-xs text-gray-600">
          <div>Time: {record.time_limit} ms</div>
          <div>Memory: {record.memory_limit} MB</div>
        </div>
      ),
    },
    {
      title: "Acceptance",
      key: "acceptance",
      width: 160,
      render: (_, record) => (
        <div className="text-xs text-gray-600">
          <div>Rate: {record.acceptance_rate.toFixed(1)}%</div>
          <div>
            {record.total_accepted}/{record.total_submissions} AC
          </div>
        </div>
      ),
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      key: "created_at",
      width: 180,
      sorter: (a, b) =>
        dayjs(a.created_at).unix() - dayjs(b.created_at).unix(),
      render: (date) => (
        <span className="font-mono text-gray-600">
          {dayjs(date).format("MMM DD, YYYY HH:mm")}
        </span>
      ),
    },
    {
      title: "Updated At",
      dataIndex: "updated_at",
      key: "updated_at",
      width: 180,
      render: (date) => (
        <span className="font-mono text-gray-600">
          {dayjs(date).format("MMM DD, YYYY HH:mm")}
        </span>
      ),
    },
    {
      title: "Actions",
      key: "actions",
      width: 170,
      render: (_, record) => (
        <div className="flex space-x-1">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleOpenModal(record)}
            disabled={deleteMutation.isPending}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete Problem"
            description={`Are you sure you want to delete "${record.title}"?`}
            icon={<ExclamationCircleOutlined style={{ color: "red" }} />}
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
            />
          </Popconfirm>
        </div>
      ),
    },
  ]

  const loading =
    isFetching ||
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending

  return (
    <div>
      <div className="flex justify-between items-center mb-6 p-4 bg-white rounded-lg shadow-sm">
        <h1 className="text-3xl font-bold">Problem Management</h1>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => handleOpenModal()}
          loading={createMutation.isPending}
          disabled={deleteMutation.isPending}
        >
          Add Problem
        </Button>
      </div>

      <div className="bg-white rounded-lg shadow-lg">
        <Table
          bordered
          rowKey="id"
          columns={columns}
          dataSource={problems}
          loading={loading}
          pagination={{
            pageSize: data?.limit || 10,
            total: data?.total,
          }}
          scroll={{ x: 1400 }}
          className="shadow-lg rounded-lg"
          size="middle"
        />
      </div>

      <Modal
        title={editingProblem ? "Edit Problem" : "Create Problem"}
        open={isModalOpen}
        onCancel={() => {
          setIsModalOpen(false)
          form.resetFields()
          setEditingProblem(null)
        }}
        footer={null}
        width={800}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          disabled={createMutation.isPending || updateMutation.isPending}
        >
          <Form.Item
            name="title"
            label="Title"
            rules={[{ required: true, message: "Please enter title" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="slug"
            label="Slug"
            rules={[{ required: true, message: "Please enter slug" }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="difficulty"
            label="Difficulty"
            rules={[{ required: true, message: "Please select difficulty" }]}
          >
            <Select
              options={[
                { label: "Easy", value: "easy" },
                { label: "Medium", value: "medium" },
                { label: "Hard", value: "hard" },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="status"
            label="Status"
            rules={[{ required: true, message: "Please select status" }]}
          >
            <Select
              options={[
                { label: "Draft", value: "draft" },
                { label: "Published", value: "published" },
              ]}
            />
          </Form.Item>

          <Form.Item
            name="time_limit"
            label="Time Limit (ms)"
            rules={[{ required: true, message: "Please enter time limit" }]}
          >
            <InputNumber min={1} className="w-full" />
          </Form.Item>

          <Form.Item
            name="memory_limit"
            label="Memory (MB)"
            rules={[{ required: true, message: "Please enter memory limit" }]}
          >
            <InputNumber min={1} className="w-full" />
          </Form.Item>

          <Form.Item
            name="validator_type"
            label="Validator Type"
            rules={[{ required: true, message: "Please select validator type" }]}
          >
            <Select
              options={[{ label: "Exact Match", value: "exact_match" }]}
            />
          </Form.Item>

          <Form.Item
            name="description"
            label="Description"
            rules={[{ required: true, message: "Please enter description" }]}
          >
            <Input.TextArea rows={4} />
          </Form.Item>

          <Form.Item
            name="input_format"
            label="Input Format"
            rules={[{ message: "Please enter input format" }]}
          >
            <Input.TextArea rows={3} />
          </Form.Item>

          <Form.Item
            name="output_format"
            label="Output Format"
            rules={[{ message: "Please enter output format" }]}
          >
            <Input.TextArea rows={3} />
          </Form.Item>

          <Form.Item
            name="constraints"
            label="Constraints"
            rules={[{ message: "Please enter constraints" }]}
          >
            <Input.TextArea rows={3} />
          </Form.Item>

          <Form.Item name="is_active" label="Active" valuePropName="checked">
            <Switch
              checkedChildren={<CheckCircleOutlined />}
              unCheckedChildren={<CloseCircleOutlined />}
            />
          </Form.Item>

          <Form.Item className="mb-0">
            <div className="flex justify-end space-x-2">
              <Button
                onClick={() => {
                  setIsModalOpen(false)
                  form.resetFields()
                  setEditingProblem(null)
                }}
              >
                Cancel
              </Button>
              <Button
                type="primary"
                htmlType="submit"
                loading={
                  createMutation.isPending || updateMutation.isPending
                }
              >
                {editingProblem ? "Update Problem" : "Create Problem"}
              </Button>
            </div>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
