import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { Modal, Form, Input, Button, Table, Switch, Popconfirm, Tag } from "antd";
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, EyeInvisibleOutlined } from "@ant-design/icons";
import { useState } from "react";
import toast from "react-hot-toast";
import type { CreateTestCaseRequest } from "../../../types/request";
import type { TestCase } from "../../../types";
import { adminTestcaseApi } from "../../../api/adminApi";

interface TestCaseFormValues extends CreateTestCaseRequest {}

export interface TestCaseListProps {
  problemId: number;
}

export default function TestCaseList({ problemId }: TestCaseListProps) {
  const queryClient = useQueryClient();
  const [form] = Form.useForm<TestCaseFormValues>();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTestCase, setEditingTestCase] = useState<TestCase | null>(null);

  const { data, isFetching } = useQuery({
    queryKey: ["testcases", problemId],
    queryFn: () => adminTestcaseApi.getAll(problemId),
  });

  // console.log(data)
  const testCases = data?.data.data || [];
  console.log(testCases)

  const createMutation = useMutation({
    mutationFn: (values: CreateTestCaseRequest) => adminTestcaseApi.create(problemId, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case created");
      closeModal();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ testcaseId, values }: { testcaseId: number; values: CreateTestCaseRequest }) =>
      adminTestcaseApi.update(problemId, testcaseId, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case updated");
      closeModal();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (testcaseId: number) => adminTestcaseApi.delete(problemId, testcaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case deleted");
    },
  });

  const openModal = (testcase?: TestCase) => {
    if (testcase) {
      form.setFieldsValue(testcase);
      setEditingTestCase(testcase);
    } else {
      form.resetFields();
      setEditingTestCase(null);
    }
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    form.resetFields();
    setEditingTestCase(null);
  };

  const handleSubmit = (values: TestCaseFormValues) => {
    values.problem_id = problemId
    if (editingTestCase) {
      updateMutation.mutate({ testcaseId: editingTestCase.id, values });
    } else {
      createMutation.mutate(values);
    }
  };

  const columns = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 80,
    },
    {
      title: "Order",
      dataIndex: "order",
      key: "order",
      width: 80,
    },
    {
      title: "Input",
      dataIndex: "input",
      key: "input",
      render: (input: string) => (
        <div className="max-w-[200px] truncate bg-gray-50 p-2 rounded font-mono text-sm">
          {input}
        </div>
      ),
    },
    {
      title: "Expected Output",
      dataIndex: "expected_output",
      key: "expected_output",
      render: (output: string) => (
        <div className="max-w-[200px] truncate bg-green-50 p-2 rounded font-mono text-sm">
          {output}
        </div>
      ),
    },
    {
      title: "Sample",
      dataIndex: "is_sample",
      key: "is_sample",
      width: 90,
      render: (isSample: boolean) => (
        <Tag color={isSample ? "blue" : "default"}>{isSample ? "SAMPLE" : "-"}</Tag>
      ),
    },
    {
      title: "Hidden",
      dataIndex: "is_hidden",
      key: "is_hidden",
      width: 90,
      render: (isHidden: boolean) => (
        <Switch
          checked={isHidden}
          checkedChildren={<EyeInvisibleOutlined />}
          unCheckedChildren={<EyeOutlined />}
          size="small"
        />
      ),
    },
    {
      title: "Actions",
      key: "actions",
      width: 120,
      render: (_, record: TestCase) => (
        <div className="flex space-x-1">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => openModal(record)}
          />
          <Popconfirm
            title="Delete Test Case"
            description="Are you sure?"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </div>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">Test Cases ({testCases.length})</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => openModal()}
          loading={createMutation.isPending}
        >
          Add Test Case
        </Button>
      </div>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={testCases}
        loading={isFetching}
        pagination={false}
        size="small"
        scroll={{ x: 1000 }}
      />

      <Modal
        title={editingTestCase ? "Edit Test Case" : "Create Test Case"}
        open={isModalOpen}
        onCancel={closeModal}
        footer={null}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="input"
            label="Input"
            rules={[{ required: true, message: "Input is required" }]}
          >
            <Input.TextArea rows={4} placeholder="Enter test case input" />
          </Form.Item>
          <Form.Item
            name="expected_output"
            label="Expected Output"
            rules={[{ required: true, message: "Expected output is required" }]}
          >
            <Input.TextArea rows={4} placeholder="Enter expected output" />
          </Form.Item>
          <Form.Item name="is_sample" label="Sample Test Case" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item name="is_hidden" label="Hidden" valuePropName="checked">
            <Switch />
          </Form.Item>
          <Form.Item className="mb-0">
            <div className="flex justify-end space-x-2">
              <Button onClick={closeModal}>Cancel</Button>
              <Button type="primary" htmlType="submit" loading={createMutation.isPending || updateMutation.isPending}>
                {editingTestCase ? "Update" : "Create"}
              </Button>
            </div>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
