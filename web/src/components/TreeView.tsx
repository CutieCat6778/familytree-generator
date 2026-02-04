import React, { useEffect, useRef, useState } from 'react';
import * as d3 from 'd3';
import { VisualizationData, VisualizationNode } from '../types';
import { calculateTreeDimensions, getPersonColor } from '../utils/treeLayout';

interface TreeViewProps {
  data: VisualizationData;
  onPersonClick: (person: VisualizationNode) => void;
  selectedPersonId: string | null;
}

export const TreeView: React.FC<TreeViewProps> = ({ data, onPersonClick, selectedPersonId }) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });
  const [svgWidth, setSvgWidth] = useState(800);

  useEffect(() => {

    const updateDimensions = () => {
      if (svgRef.current?.parentElement) {
        setDimensions({
          width: svgRef.current.parentElement.clientWidth,
          height: svgRef.current.parentElement.clientHeight,
        });
      }
    };

    updateDimensions();
    window.addEventListener('resize', updateDimensions);
    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  useEffect(() => {
    if (!svgRef.current || !data) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    const { width, height } = dimensions;
    const margin = { top: 40, right: 40, bottom: 40, left: 40 };


    const nodeMap = new Map<string, VisualizationNode>();
    data.nodes.forEach(n => nodeMap.set(n.id, n));


    const childrenMap = new Map<string, string[]>();
    data.edges.filter(e => e.type === 'parent').forEach(e => {
      const children = childrenMap.get(e.source) || [];
      children.push(e.target);
      childrenMap.set(e.source, children);
    });


    interface TreeNodeData extends VisualizationNode {
      children?: TreeNodeData[];
    }

    function buildTree(nodeId: string, visited: Set<string>): TreeNodeData | null {
      if (visited.has(nodeId)) return null;
      visited.add(nodeId);

      const node = nodeMap.get(nodeId);
      if (!node) return null;

      const treeNode: TreeNodeData = { ...node, children: [] };
      const childIds = childrenMap.get(nodeId) || [];

      childIds.forEach(childId => {
        const child = buildTree(childId, visited);
        if (child) {
          treeNode.children!.push(child);
        }
      });

      if (treeNode.children!.length === 0) {
        delete treeNode.children;
      }

      return treeNode;
    }

    const rootData = buildTree(data.root_id, new Set());
    if (!rootData) return;


    const root = d3.hierarchy(rootData);
    const leafNodes = root.leaves();
    const leafCount = leafNodes.length;
    const leafIds = new Set(leafNodes.map(node => node.data.id));
    const mainNodeIds = new Set(root.descendants().map(node => node.data.id));
    const extraLeafSpouseIds = new Set<string>();

    data.edges.filter(e => e.type === 'spouse').forEach(edge => {
      const sourceIn = mainNodeIds.has(edge.source);
      const targetIn = mainNodeIds.has(edge.target);
      if (sourceIn && !targetIn && leafIds.has(edge.source)) {
        extraLeafSpouseIds.add(edge.target);
      } else if (!sourceIn && targetIn && leafIds.has(edge.target)) {
        extraLeafSpouseIds.add(edge.source);
      }
    });

    const minLayout = calculateTreeDimensions(data.nodes.length);
    const minLeafSpacing = extraLeafSpouseIds.size > 0 ? 280 : 160;
    const minLeafRowWidth = Math.max(1, leafCount - 1) * minLeafSpacing;
    const minSvgWidth = Math.max(minLayout.width, minLeafRowWidth + margin.left + margin.right);
    const layoutWidth = Math.max(width, minSvgWidth);
    const layoutHeight = height;

    setSvgWidth((current) => (Math.abs(current - layoutWidth) < 1 ? current : layoutWidth));


    const treeLayout = d3.tree<TreeNodeData>()
      .size([layoutWidth - margin.left - margin.right, layoutHeight - margin.top - margin.bottom])
      .separation((a, b) => (a.parent === b.parent ? 1.5 : 2));

    const treeData = treeLayout(root);
    const nodeById = new Map<string, d3.HierarchyPointNode<TreeNodeData>>(
      treeData.descendants().map(node => [node.data.id, node]),
    );


    const zoom = d3.zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.3, 3])
      .on('zoom', (event) => {
        g.attr('transform', event.transform);
      });

    svg.call(zoom);


    const g = svg.append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);


    g.selectAll('.link')
      .data(treeData.links())
      .join('path')
      .attr('class', 'link')
      .attr('fill', 'none')
      .attr('stroke', '#999')
      .attr('stroke-width', 2)
      .attr('d', d3.linkVertical<d3.HierarchyPointLink<TreeNodeData>, d3.HierarchyPointNode<TreeNodeData>>()
        .x(d => d.x)
        .y(d => d.y)
      );


    const spouseEdges = data.edges.filter(e => e.type === 'spouse');
    const extraSpouseNodes: Array<{
      data: TreeNodeData;
      x: number;
      y: number;
    }> = [];
    const extraSpouseIds = new Set<string>();
    const spouseOffset = 140;

    spouseEdges.forEach(edge => {
      const source = nodeById.get(edge.source);
      const target = nodeById.get(edge.target);

      if (source && !target) {
        if (!extraSpouseIds.has(edge.target)) {
          const spouseData = nodeMap.get(edge.target);
          if (spouseData) {
            extraSpouseNodes.push({
              data: { ...spouseData },
              x: source.x + spouseOffset,
              y: source.y,
            });
            extraSpouseIds.add(edge.target);
          }
        }
      } else if (!source && target) {
        if (!extraSpouseIds.has(edge.source)) {
          const spouseData = nodeMap.get(edge.source);
          if (spouseData) {
            extraSpouseNodes.push({
              data: { ...spouseData },
              x: target.x - spouseOffset,
              y: target.y,
            });
            extraSpouseIds.add(edge.source);
          }
        }
      }
    });

    const extraSpousePos = new Map<string, { x: number; y: number }>(
      extraSpouseNodes.map(node => [node.data.id, { x: node.x, y: node.y }]),
    );

    const getNodePosition = (id: string) => {
      const mainNode = nodeById.get(id);
      if (mainNode) {
        return { x: mainNode.x, y: mainNode.y };
      }
      return extraSpousePos.get(id) ?? null;
    };

    spouseEdges.forEach(edge => {
      const source = getNodePosition(edge.source);
      const target = getNodePosition(edge.target);

      if (source && target) {
        g.append('line')
          .attr('class', 'spouse-link')
          .attr('x1', source.x)
          .attr('y1', source.y)
          .attr('x2', target.x)
          .attr('y2', target.y)
          .attr('stroke', '#ff69b4')
          .attr('stroke-width', 2)
          .attr('stroke-dasharray', '5,5');
      }
    });


    const nodes = g.selectAll('.node')
      .data(treeData.descendants())
      .join('g')
      .attr('class', 'node')
      .attr('transform', d => `translate(${d.x},${d.y})`)
      .style('cursor', 'pointer')
      .on('click', (_, d) => {
        onPersonClick(d.data);
      });


    nodes.append('rect')
      .attr('x', -60)
      .attr('y', -25)
      .attr('width', 120)
      .attr('height', 50)
      .attr('rx', 8)
      .attr('fill', d => getPersonColor(d.data.gender, d.data.is_alive))
      .attr('stroke', d => d.data.id === selectedPersonId ? '#333' : d.data.id === data.root_id ? 'gold' : 'transparent')
      .attr('stroke-width', d => (d.data.id === selectedPersonId || d.data.id === data.root_id) ? 3 : 0);


    nodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', -5)
      .attr('fill', '#fff')
      .attr('font-size', '12px')
      .attr('font-weight', 'bold')
      .text(d => {
        const name = d.data.name;
        return name.length > 15 ? name.substring(0, 12) + '...' : name;
      });


    nodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', 12)
      .attr('fill', 'rgba(255,255,255,0.9)')
      .attr('font-size', '11px')
      .text(d => {
        return d.data.death_year
          ? `${d.data.birth_year}-${d.data.death_year}`
          : `b. ${d.data.birth_year}`;
      });


    const spouseNodes = g.selectAll('.spouse-node')
      .data(extraSpouseNodes)
      .join('g')
      .attr('class', 'node spouse-node')
      .attr('transform', d => `translate(${d.x},${d.y})`)
      .style('cursor', 'pointer')
      .on('click', (_, d) => {
        onPersonClick(d.data);
      });

    spouseNodes.append('rect')
      .attr('x', -60)
      .attr('y', -25)
      .attr('width', 120)
      .attr('height', 50)
      .attr('rx', 8)
      .attr('fill', d => getPersonColor(d.data.gender, d.data.is_alive))
      .attr('stroke', d => d.data.id === selectedPersonId ? '#333' : 'transparent')
      .attr('stroke-width', d => (d.data.id === selectedPersonId) ? 3 : 0);

    spouseNodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', -5)
      .attr('fill', '#fff')
      .attr('font-size', '12px')
      .attr('font-weight', 'bold')
      .text(d => {
        const name = d.data.name;
        return name.length > 15 ? name.substring(0, 12) + '...' : name;
      });

    spouseNodes.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', 12)
      .attr('fill', 'rgba(255,255,255,0.9)')
      .attr('font-size', '11px')
      .text(d => {
        return d.data.death_year
          ? `${d.data.birth_year}-${d.data.death_year}`
          : `b. ${d.data.birth_year}`;
      });


    const bounds = g.node()?.getBBox();
    if (bounds) {
      const scale = 1;
      const translateX = margin.left - bounds.x * scale;
      const translateY = margin.top - bounds.y * scale;
      svg.call(zoom.transform, d3.zoomIdentity.translate(translateX, translateY).scale(scale));
    }

  }, [data, dimensions, onPersonClick, selectedPersonId]);

  return (
    <svg
      ref={svgRef}
      width={svgWidth}
      height={dimensions.height}
      style={{ background: '#fafafa', borderRadius: '8px' }}
    />
  );
};
