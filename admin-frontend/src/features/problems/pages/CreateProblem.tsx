import { useState, useEffect } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import ProblemForm from "../components/ProblemForm"
import ProblemPreview from "./ProblemPreview"
import type { CreateOrUpdateProblemRequest } from "../../../types/request"
import { adminProblemApi } from "../../../api/adminApi"
import toast from "react-hot-toast"

export default function CreateProblem() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const isEditMode = Boolean(id)

  const [formData, setFormData] = useState<CreateOrUpdateProblemRequest>({
    title: "",
    slug: "",
    description: "",
    difficulty: "easy",
    time_limit: 1000,
    memory_limit: 256,
    validator_type: "exact_match",
    input_format: "",
    output_format: "",
    constraints: "",
    status: "draft",
    is_active: false,
  })

  // Fetch problem data if in edit mode
  const { data: problemData, isLoading: isFetchingProblem } = useQuery({
    queryKey: ['problem', id],
    queryFn: () => adminProblemApi.getById(parseInt(id!)),
    enabled: isEditMode && !!id,
    select: (data) => data.data.data,
  })

  // Set form data when problem is fetched
  useEffect(() => {
    if (problemData) {
      setFormData(problemData)
    }
  }, [problemData])

  // Create mutation
  const createMutation = useMutation({
    mutationFn: (data: CreateOrUpdateProblemRequest) => 
      adminProblemApi.create(data),
    onSuccess: (response) => {
      const problemId = response.data.id
      toast.success("Problem created successfully!")
      queryClient.invalidateQueries({ queryKey: ['problems'] })
      navigate(`/problems/${problemId}/testcases`)
    },
    onError: (error) => {
      console.error("Error creating problem:", error)
      toast.error("Failed to create problem")
    },
  })

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: (data: CreateOrUpdateProblemRequest) => 
      adminProblemApi.update(parseInt(id!), data),
    onSuccess: (response) => {
      const problemId = response.data.id || id
      toast.success("Problem updated successfully!")
      queryClient.invalidateQueries({ queryKey: ['problem', id] })
      queryClient.invalidateQueries({ queryKey: ['problems'] })
      navigate(`/problems/${problemId}/testcases`)
    },
    onError: (error) => {
      console.error("Error updating problem:", error)
      toast.error("Failed to update problem")
    },
  })

  // Save draft mutation (doesn't redirect)
  const saveDraftMutation = useMutation({
    mutationFn: (data: CreateOrUpdateProblemRequest) => {
      const draftData = { ...data, status: "draft" as const }
      return isEditMode && id
        ? adminProblemApi.update(parseInt(id), draftData)
        : adminProblemApi.create(draftData)
    },
    onSuccess: (response) => {
      const problemId = response.data.id || id
      toast.success("Draft saved successfully!")
      queryClient.invalidateQueries({ queryKey: ['problems'] })
      
      // If creating new draft, redirect to edit mode
      if (!isEditMode) {
        navigate(`/admin/problems/edit/${problemId}`, { replace: true })
      }
    },
    onError: (error) => {
      console.error("Error saving draft:", error)
      toast.error("Failed to save draft")
    },
  })

  const handleFormChange = (updates: Partial<CreateOrUpdateProblemRequest>) => {
    setFormData(prev => ({ ...prev, ...updates }))
  }

  const handleSubmit = () => {
    if (isEditMode) {
      updateMutation.mutate(formData)
    } else {
      createMutation.mutate(formData)
    }
  }

  const handleSaveDraft = () => {
    saveDraftMutation.mutate(formData)
  }

  const isLoading = 
    createMutation.isPending || 
    updateMutation.isPending || 
    saveDraftMutation.isPending

  if (isFetchingProblem) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-lg">Loading problem...</div>
      </div>
    )
  }

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <div className="bg-white border-b px-6 py-4">
        <h1 className="text-2xl font-bold">
          {isEditMode ? "Edit Problem" : "Create Problem"}
        </h1>
      </div>

      {/* Main Content */}
      <div className="flex gap-4 flex-1 overflow-hidden">
        <div className="flex-1 p-4 overflow-y-auto">
          <ProblemForm 
            formData={formData} 
            onChange={handleFormChange}
            onSubmit={handleSubmit}
            onSaveDraft={handleSaveDraft}
            loading={isLoading}
            isEditMode={isEditMode}
          />
        </div>
        <div className="flex-1 p-4 overflow-y-auto border-l bg-gray-50">
          <ProblemPreview formData={formData} />
        </div>
      </div>
    </div>
  )
}
