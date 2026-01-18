import TiptapEditor from "../../../components/editor/TiptapEditor"
import type { CreateOrUpdateProblemRequest } from "../../../types/request"
import { useQuery } from "@tanstack/react-query"
import { adminProblemApi } from "../../../lib/api/admin"
import type { Tag, Category } from "../../../types"

interface ProblemFormProps {
  formData: CreateOrUpdateProblemRequest
  onChange: (updates: Partial<CreateOrUpdateProblemRequest>) => void
  onSubmit: () => void
  onSaveDraft: () => void
  loading: boolean
  isEditMode: boolean
}

export default function ProblemForm({
  formData,
  onChange,
  onSubmit,
  onSaveDraft,
  loading,
  isEditMode
}: ProblemFormProps) {
  const { data: tagsResponse } = useQuery({
    queryKey: ['tags'],
    queryFn: () => adminProblemApi.getTags(),
  })

  const { data: categoriesResponse } = useQuery({
    queryKey: ['categories'],
    queryFn: () => adminProblemApi.getCategories(),
  })

  const tags = tagsResponse?.data?.data || []
  const categories = categoriesResponse?.data?.data || []

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit()
  }

  return (
    <div className="space-y-6">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Title */}
        <div>
          <label htmlFor="title" className="block text-sm font-medium mb-2">
            Title <span className="text-red-500">*</span>
          </label>
          <input
            id="title"
            type="text"
            value={formData.title}
            onChange={(e) => onChange({ title: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Two Sum"
            required
            disabled={loading}
          />
        </div>

        {/* Slug */}
        <div>
          <label htmlFor="slug" className="block text-sm font-medium mb-2">
            Slug <span className="text-red-500">*</span>
          </label>
          <input
            id="slug"
            type="text"
            value={formData.slug}
            onChange={(e) => onChange({ slug: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="two-sum"
            required
            disabled={loading}
          />
        </div>

        {/* Difficulty */}
        <div>
          <label htmlFor="difficulty" className="block text-sm font-medium mb-2">
            Difficulty <span className="text-red-500">*</span>
          </label>
          <select
            id="difficulty"
            value={formData.difficulty}
            onChange={(e) => onChange({ difficulty: e.target.value as "easy" | "medium" | "hard" })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={loading}
          >
            <option value="easy">Easy</option>
            <option value="medium">Medium</option>
            <option value="hard">Hard</option>
          </select>
        </div>

        {/* Description */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Description <span className="text-red-500">*</span>
          </label>
          <TiptapEditor
            content={formData.description}
            onChange={(content) => onChange({ description: content })}
          />
        </div>

        {/* Input Format */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Input Format
          </label>
          <TiptapEditor
            content={formData.input_format}
            onChange={(content) => onChange({ input_format: content })}
          />
        </div>

        {/* Output Format */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Output Format
          </label>
          <TiptapEditor
            content={formData.output_format}
            onChange={(content) => onChange({ output_format: content })}
          />
        </div>

        {/* Constraints */}
        <div>
          <label className="block text-sm font-medium mb-2">
            Constraints
          </label>
          <TiptapEditor
            content={formData.constraints}
            onChange={(content) => onChange({ constraints: content })}
          />
        </div>

        {/* Time and Memory Limits */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="time_limit" className="block text-sm font-medium mb-2">
              Time Limit (ms) <span className="text-red-500">*</span>
            </label>
            <input
              id="time_limit"
              type="number"
              value={formData.time_limit}
              onChange={(e) => onChange({ time_limit: parseInt(e.target.value) })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              min="100"
              disabled={loading}
            />
          </div>
          <div>
            <label htmlFor="memory_limit" className="block text-sm font-medium mb-2">
              Memory Limit (MB) <span className="text-red-500">*</span>
            </label>
            <input
              id="memory_limit"
              type="number"
              value={formData.memory_limit}
              onChange={(e) => onChange({ memory_limit: parseInt(e.target.value) })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
              min="64"
              disabled={loading}
            />
          </div>
        </div>

        {/* Status and Active */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="status" className="block text-sm font-medium mb-2">
              Status
            </label>
            <select
              id="status"
              value={formData.status}
              onChange={(e) => onChange({ status: e.target.value as "draft" | "published" })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            >
              <option value="draft">Draft</option>
              <option value="published">Published</option>
            </select>
          </div>
          <div className="flex items-center pt-8">
            <input
              id="is_active"
              type="checkbox"
              checked={formData.is_active}
              onChange={(e) => onChange({ is_active: e.target.checked })}
              className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              disabled={loading}
            />
            <label htmlFor="is_active" className="ml-2 text-sm font-medium">
              Is Active
            </label>
          </div>
        </div>

        {/* Tags */}
        <div className="border-t pt-6">
          <label className="block text-sm font-medium mb-3">
            Tags (Topics)
          </label>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3 p-4 bg-gray-50 rounded-lg max-h-60 overflow-y-auto border">
            {tags.map((tag: Tag) => (
              <label key={tag.id} className="flex items-center space-x-2 cursor-pointer hover:bg-white p-1 rounded transition-colors">
                <input
                  type="checkbox"
                  checked={formData.tag_ids?.includes(tag.id)}
                  onChange={(e) => {
                    const currentIds = formData.tag_ids || []
                    const nextIds = e.target.checked
                      ? [...currentIds, tag.id]
                      : currentIds.filter(id => id !== tag.id)
                    onChange({ tag_ids: nextIds })
                  }}
                  className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  disabled={loading}
                />
                <span className="text-sm text-gray-700">{tag.name}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Categories */}
        <div className="border-t pt-6">
          <label className="block text-sm font-medium mb-3">
            Categories
          </label>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-3 p-4 bg-gray-50 rounded-lg border">
            {categories.map((cat: Category) => (
              <label key={cat.id} className="flex items-center space-x-2 cursor-pointer hover:bg-white p-1 rounded transition-colors">
                <input
                  type="checkbox"
                  checked={formData.category_ids?.includes(cat.id)}
                  onChange={(e) => {
                    const currentIds = formData.category_ids || []
                    const nextIds = e.target.checked
                      ? [...currentIds, cat.id]
                      : currentIds.filter(id => id !== cat.id)
                    onChange({ category_ids: nextIds })
                  }}
                  className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  disabled={loading}
                />
                <span className="text-sm text-gray-700">{cat.name}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Submit Buttons */}
        <div className="flex gap-4 pt-4 border-t">
          <button
            type="submit"
            disabled={loading}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Saving..." : isEditMode ? "Update Problem" : "Create Problem"}
          </button>
          <button
            type="button"
            onClick={onSaveDraft}
            disabled={loading}
            className="px-6 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Saving..." : "Save as Draft"}
          </button>
        </div>
      </form>
    </div>
  )
}
